package script

import (
	"fmt"
	"os"
	"errors"
)

var (
	ErrEarlyEOF = errors.New("Unexpected EOF when reading OP arguments")
	ErrInvalidInstruction = errors.New("Invalid instruction")
	ErrNavigation = errors.New("SmartParse navigation error")
)

type Parser struct {
	rawinput []byte
	current int
	startAddr int

	script *Script
	cdl    *CodeDataLog
}

func ParseFile(filename string, startAddr int, cdl *CodeDataLog) (*Script, error) {
	rawfile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return Parse(rawfile, startAddr, cdl)
}

func SmartParseFile(filename string, startAddr int, cdl *CodeDataLog) (*Script, error) {
	rawfile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return SmartParse(rawfile, startAddr, cdl)
}

func SmartParse(rawinput []byte, startAddr int, cdl *CodeDataLog) (*Script, error) {
	if len(rawinput) < 3 {
		return nil, fmt.Errorf("not enough bytes for script")
	}

	p := &Parser{

		script: &Script{
			Tokens: []*Token{},
			Warnings: []string{},
			StackAddress: (int(rawinput[1])<<8) | int(rawinput[0]),
			StartAddress: startAddr,
			Labels: make(map[int]*Label), // map[location]name

			CDL: cdl,
			origSize: len(rawinput),
		},

		rawinput: rawinput,
		startAddr: startAddr,
	}

	if p.script.CDL == nil {
		p.script.CDL = NewCDL()
	}

	tokenMap := make(map[int]*Token)

	// starting point is the third byte in the script.
	branches := []int{ 2 }
	for _, ent := range p.script.CDL.getEntries() {
		addr := ent-startAddr
		if addr > 0 {
			branches = append(branches, addr)
		}
	}


	visited := make([]bool, len(p.rawinput))

	for len(branches) > 0 {
		st := branches[0]+startAddr
		//fmt.Printf("start @ $%04X\n", st)
		p.script.Labels[st] = AutoLabel(st)

INNER:
		for p.current = branches[0]; p.current < len(p.rawinput); p.current++ {
			//branches = branches[1:]

			if p.current < 0 {
				return p.script, errors.Join(ErrNavigation,
					fmt.Errorf("HOW IS CURRENT NEGATIVE?????"))
			}

			if visited[p.current] {
				//fmt.Printf("found visited at $%04X\n", p.current+startAddr)
				break
			}

			visited[p.current] = true
			raw := p.rawinput[p.current]

			token := &Token{
				Offset: startAddr+p.current,
				Raw: raw,
				Inline: []InlineVal{},
			}
			p.script.Tokens = append(p.script.Tokens, token)
			tokenMap[token.Offset] = token

			//fmt.Printf("{$%04X} %s\n", token.Offset, token.String(map[int]*Label{}))

			p.script.CDL.setCode(p.current+p.startAddr)
			if raw < 0x80 {
				continue
			}

			err := p.parseToken(token, raw)
			if err != nil {
				return p.script, err
			}

			//fmt.Println(token.String(map[int]*Label{}))

			switch raw {
			case 0x86, 0xAC, 0xAA, 0xFF, 0x81, 0x9B, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xFD: // return, long_return, long_jump, break_engine & halts
				//fmt.Printf("[$%04X] %s\n",
				//	token.Offset, token.Instruction.Name)
				break INNER

			case 0x84, 0xBF, 0xC0, 0x85: // jump_abs, jump_not_zero, jump_zero, call_abs
				if len(token.Inline) < 1 {
					return p.script, errors.Join(ErrNavigation,
						fmt.Errorf("jump missing target"))
				}

				if len(token.Inline) > 1 {
					return p.script, errors.Join(ErrNavigation,
						fmt.Errorf("jump has too many targets"))
				}

				val := token.Inline[0].Int()
				//fmt.Printf("[$%04X] %s $%04X\n",
				//	token.Offset, token.Instruction.Name, val)
				branches = append(branches, val-startAddr)
				p.script.Labels[val] = AutoLabel(val)

				if raw == 0x84 { // not jump_abs
					break INNER
				}

			case 0xC1, 0xEE: // jump_switch, call_switch
				if len(token.Inline) < 2 {
					return p.script, errors.Join(ErrNavigation,
						fmt.Errorf("switch missing targets"))
				}

				count := token.Inline[0].Int()
				if len(token.Inline) != count+1 {
					return p.script, errors.Join(ErrNavigation,
						fmt.Errorf("switch target missmatch (expected %d, got %d)", count, len(token.Inline)-1))
				}

				for _, val := range token.Inline[1:] {
					//fmt.Printf("[$%04X] %s $%04X\n",
					//	token.Offset, token.Instruction.Name, val.Int())
					branches = append(branches, val.Int()-startAddr)
					p.script.Labels[val.Int()] = AutoLabel(val.Int())
				}

				if raw == 0xC1 { // jump_switch
					break INNER
				}
			}

			if token.Instruction.OpCount == 2 && !token.Instruction.InlineImmediate {
				val := token.Inline[0].Int()
				if _, ok := p.script.Labels[val]; !ok {//&& val >= startAddr {
					p.script.Labels[val] = AutoLabelVar(val)
				}
				p.script.CDL.setData(val)
				p.script.CDL.setData(val+1)
			}
		}

		if len(branches) == 1 {
			break
		}
		branches = branches[1:]
	}

	// Add data tokens
	for addr, bit := range p.script.CDL.cache {
		if addr < 0x6002 {
			continue
		}

		// ignore code bytes
		if bit & cdlCode == cdlCode {
			continue
		}

		// ignore labels outside the script's address range
		if addr > len(rawinput)+0x6000 {
			continue
		}

		if _, ok := p.script.Labels[addr]; ok {
			p.script.Tokens = append(p.script.Tokens, &Token{
				Offset: addr,
				Inline: []InlineVal{NewWordVal([]byte{rawinput[addr-0x6000], rawinput[addr+1-0x6000]})},
				IsVariable: true,
				IsData: true,
				cdl: bit.String(),
			})
		} else {
			p.script.Tokens = append(p.script.Tokens, &Token{
				Offset: addr,
				Raw: rawinput[addr-0x6000],
				IsData: true,
				cdl: bit.String(),
			})
		}
	}

	return p.script, nil
}

func Parse(rawinput []byte, startAddr int, cdl *CodeDataLog) (*Script, error) {
	if len(rawinput) < 3 {
		return nil, fmt.Errorf("not enough bytes for script")
	}

	p := &Parser{
		script: &Script{
			Tokens: []*Token{},
			Warnings: []string{},
			StackAddress: (int(rawinput[1])<<8) | int(rawinput[0]),
			StartAddress: startAddr,
			Labels: make(map[int]*Label), // map[location]name
			CDL: cdl,
			origSize: len(rawinput),
		},
		rawinput: rawinput,
		startAddr: startAddr,
	}
	tokenMap := make(map[int]*Token)

	if p.script.CDL == nil {
		p.script.CDL = NewCDL()
	}

	//earliestVar := len(p.rawinput)-2
	//fmt.Printf("var start bounds: $%04X, $%04X\n", startAddr, startAddr+len(p.rawinput))

	for p.current = 2; p.current < len(p.rawinput); p.current++ {
		//if p.current >= earliestVar {
		//	fmt.Printf("Earliest Variable found at offset %d ($%04X)\n", p.current, startAddr+p.current)
		//	break
		//}

		raw := p.rawinput[p.current]

		token := &Token{
			Offset: startAddr+p.current,
			Raw: raw,
			Inline: []InlineVal{},
		}
		p.script.Tokens = append(p.script.Tokens, token)
		tokenMap[token.Offset] = token

		if raw < 0x80 || p.script.CDL.IsData(p.current+startAddr) { // || p.current >= earliestVar {
			if p.script.CDL.IsData(p.current+startAddr) {
				token.IsData = true
				//fmt.Print(".")
				//fmt.Printf("%#v\n", token)
			}
			continue
		}

		err := p.parseToken(token, raw)
		if err != nil {
			return p.script, err
		}
	}

	// Find and mark labels for a few instructions
	for _, t := range p.script.Tokens {
		if t.Instruction == nil {
			continue
		}

		switch t.Raw {
		case 0x84, 0x85, 0xBF, 0xC0: // jmp/call
			if len(t.Inline) == 0 {
				//return nil, fmt.Errorf("jump/call missing address ($%04X)", t.Offset)
				p.script.Warnings = append(p.script.Warnings,
					fmt.Sprintf("jump/call missing addresses ($%04X)", t.Offset))
				continue
			}

			addr := t.Inline[0].Int()
			found := false
			for _, tok := range p.script.Tokens {
				if tok.Offset == addr {
					tok.IsTarget = true
					found = true
					p.script.Labels[addr] = AutoLabel(addr) //fmt.Sprintf("L%04X", addr)
					break
				}
			}

			if !found {
				p.script.Warnings = append(p.script.Warnings, fmt.Sprintf("Warning: no target found for jump/call at offset $%04X; value $%04X", t.Offset, addr))
			}

		case 0xC1, 0xEE: // switches
			if len(t.Inline) < 2 {
				//return nil, fmt.Errorf("jump/call switch missing addresses")
				p.script.Warnings = append(p.script.Warnings,
					fmt.Sprintf("jump/call switch missing addresses ($%04X)", t.Offset))
				continue
			}

			for _, v := range t.Inline[1:] {
				addr := v.Int()
				found := false
				for _, tok := range p.script.Tokens {
					if tok.Offset == addr {
						tok.IsTarget = true
						found = true
						p.script.Labels[addr] = AutoLabel(addr) //fmt.Sprintf("L%04X", addr)
						break
					}
				}

				if !found {
					p.script.Warnings = append(p.script.Warnings, fmt.Sprintf("Warning: no target found for jump/call switch at offset $%04X; value: $%04X", t.Offset, addr))
				}
			}

		default:
			// if word arg, see if it's something in this script
			if t.Instruction == nil {
				//if t.IsData {
				//	fmt.Print(",")
				//}
				continue
			}

			if t.Instruction.OpCount == 2 && !t.Instruction.InlineImmediate {
				addr := t.Inline[0].Int()
				if _, ok := tokenMap[addr]; ok {
					//tok.IsVariable = true
					p.script.Labels[addr] = AutoLabelVar(addr) //fmt.Sprintf("Var_%04X", addr)
				}
			}
		}
	}

	return p.script, nil
}

func (p *Parser) parseToken(token *Token, raw byte) error {
	op, ok := InstrMap[raw]
	if !ok {
		return errors.Join(ErrInvalidInstruction,
			fmt.Errorf("OP 0x%02X not in instruction map", raw))
	}
	token.Instruction = op

	args := []InlineVal{}
	switch op.OpCount {
	case -1: // null terminated
		for ; p.current < len(p.rawinput); p.current++ {
			p.script.CDL.setCode(p.current+p.startAddr)
			val := ByteVal(p.rawinput[p.current])
			args = append(args, val)
			if p.rawinput[p.current] == 0x00 {
				break
			}
		}

	case -2: // count then count words
		// FIXME: wtf makes this different from -3??
		p.current++

		l :=  int(p.rawinput[p.current])
		p.script.CDL.setCode(p.current+p.startAddr)
		args = append(args, ByteVal(l))
		p.current++

		for c := 0; c < l; c++ {
			if len(p.rawinput) <= p.current+1 {
				return errors.Join(ErrEarlyEOF,
					fmt.Errorf("OP early end at offset 0x%X (%d) {%d} %#v", p.current, p.current, l, op))
			}

			args = append(args, WordVal([2]byte{p.rawinput[p.current], p.rawinput[p.current+1]}))
			p.script.CDL.setCode(p.current+p.startAddr)
			p.script.CDL.setCode(p.current+p.startAddr+1)
			p.current+=2
		}
		p.current--

	case -3: // count then count words.  "default" is no call (skip Code_Pointer to after args)
		p.current++

		l :=  int(p.rawinput[p.current])
		args = append(args, ByteVal(l))
		p.script.CDL.setCode(p.current+p.startAddr)
		p.current++

		for c := 0; c < l; c++ {
			if len(p.rawinput) <= p.current+1 {
				return errors.Join(ErrEarlyEOF,
					fmt.Errorf("OP early end at offset 0x%X (%d) {%d} %#v", p.current, p.current, l, op))
			}

			args = append(args, WordVal([2]byte{p.rawinput[p.current], p.rawinput[p.current+1]}))
			p.script.CDL.setCode(p.current+p.startAddr)
			p.script.CDL.setCode(p.current+p.startAddr+1)
			p.current+=2
		}
		p.current--

	case 2:
		args = append(args, WordVal([2]byte{p.rawinput[p.current+1], p.rawinput[p.current+2]}))
		p.script.CDL.setCode(p.current+p.startAddr+1)
		p.script.CDL.setCode(p.current+p.startAddr+2)
		p.current+=2

		//fmt.Printf("var at $%04X\n", val.Int())
		//if val.Int() > p.startAddr && val.Int() < p.startAddr+len(p.rawinput) && p.earliestVar > val.Int() {
		//	fmt.Printf("new earliest: $%04X\n", val.Int())
		//	p.earliestVar = val.Int()
		//}

	case 1:
		p.current++
		p.script.CDL.setCode(p.current+p.startAddr)
		args = append(args, ByteVal(p.rawinput[p.current]))
	}

	token.Inline = args
	return nil
}
