package script

import (
	"fmt"
	"os"
)

type Label struct {
	Address int
	Name string
	Comment string
	FarLabel bool
}

func AutoLabel(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("L%04X", address),
	}
}

func AutoLabelVar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("Var_%04X", address),
	}
}

func AutoLabelFar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("F%04X", address),
		FarLabel: true,
	}
}

func NewLabel(address int, name string) *Label {
	return &Label{
		Address: address,
		Name: name,
	}
}

func ParseFile(filename string, startAddr int) (*Script, error) {
	rawfile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return Parse(rawfile, startAddr)
}

func Parse(rawinput []byte, startAddr int) (*Script, error) {
	if len(rawinput) < 3 {
		return nil, fmt.Errorf("not enough bytes for script")
	}

	script := &Script{
		Tokens: []*Token{},
		Warnings: []string{},
		StackAddress: (int(rawinput[1])<<8) | int(rawinput[0]),
		StartAddress: startAddr,
		Labels: make(map[int]*Label), // map[location]name
	}
	tokenMap := make(map[int]*Token)

	for i := 2; i < len(rawinput); i++ {
		raw := rawinput[i]

		token := &Token{
			Offset: startAddr+i,
			Raw: raw,
			Inline: []InlineVal{},
		}
		script.Tokens = append(script.Tokens, token)
		tokenMap[token.Offset] = token

		if raw < 0x80 {
			continue
		}

		op, ok := InstrMap[raw]
		if !ok {
			return nil, fmt.Errorf("OP 0x%02X not in instruction map", raw)
		}
		token.Instruction = op

		args := []InlineVal{}
		switch op.OpCount {
		case -1: // null terminated
			for ; i < len(rawinput); i++ {
				val := ByteVal(rawinput[i])
				args = append(args, val)
				if rawinput[i] == 0x00 {
					break
				}
			}

		case -2: // count then count words
			i++
			l :=  int(rawinput[i])
			args = append(args, ByteVal(l))
			i++
			for c := 0; c < l; c++ {
				if len(rawinput) <= i+1 {
					return script, fmt.Errorf("OP early end at offset 0x%X (%d) {%d} %#v", i, i, l, op)
				}

				args = append(args, WordVal([2]byte{rawinput[i], rawinput[i+1]}))
				i+=2
			}
			i--

		case -3: // count then count words.  "default" is no call (skip Code_Pointer to after args)
			i++
			l :=  int(rawinput[i])
			args = append(args, ByteVal(l))
			i++
			for c := 0; c < l; c++ {
				args = append(args, WordVal([2]byte{rawinput[i], rawinput[i+1]}))
				i+=2
			}
			i--

		case 2:
			args = append(args, WordVal([2]byte{rawinput[i+1], rawinput[i+2]}))
			i+=2

		case 1:
			i++
			args = append(args, ByteVal(rawinput[i]))
		}

		token.Inline = args
	}

	// Find and mark labels for a few instructions
	for _, t := range script.Tokens {
		switch t.Raw {
		case 0x84, 0x85, 0xBF, 0xC0: // jmp/call
			if len(t.Inline) == 0 {
				return nil, fmt.Errorf("jump/call missing address")
			}

			addr := t.Inline[0].Int()
			found := false
			for _, tok := range script.Tokens {
				if tok.Offset == addr {
					tok.IsTarget = true
					found = true
					script.Labels[addr] = AutoLabel(addr) //fmt.Sprintf("L%04X", addr)
					break
				}
			}

			if !found {
				script.Warnings = append(script.Warnings, fmt.Sprintf("Warning: no target found for jump/call at offset $%04X; value $%04X", t.Offset, addr))
			}

		case 0xC1, 0xEE: // switches
			if len(t.Inline) < 2 {
				return nil, fmt.Errorf("jump/call switch missing addresses")
			}

			for _, v := range t.Inline[1:] {
				addr := v.Int()
				found := false
				for _, tok := range script.Tokens {
					if tok.Offset == addr {
						tok.IsTarget = true
						found = true
						script.Labels[addr] = AutoLabel(addr) //fmt.Sprintf("L%04X", addr)
						break
					}
				}

				if !found {
					script.Warnings = append(script.Warnings, fmt.Sprintf("Warning: no target found for jump/call switch at offset $%04X; value: $%04X", t.Offset, addr))
				}
			}

		default:
			// if word arg, see if it's something in this script
			if t.Instruction == nil {
				continue
			}
			if t.Instruction.OpCount == 2 {
				addr := t.Inline[0].Int()
				if tok, ok := tokenMap[addr]; ok {
					tok.IsVariable = true
					script.Labels[addr] = AutoLabelVar(addr) //fmt.Sprintf("Var_%04X", addr)
				}
			}
		}
	}

	return script, nil
}
