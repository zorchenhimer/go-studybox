package script

import (
	"fmt"
	"strings"
)

type Token struct {
	Offset   int // in CPU space
	Raw      byte
	Inline   []InlineVal
	IsTarget bool   // target of a call/jump?
	IsVariable bool // target of something else
	IsData     bool // from CDL

	Instruction *Instruction
}

func (t Token) String(labels map[int]*Label) string {
	suffix := ""
	switch t.Raw {
	case 0x86: // Newline after return
		suffix = "\n"
	}

	prefix := ""
	if lbl, ok := labels[t.Offset]; ok {
		comment := ""
		if lbl.Comment != "" {
			comment = "; "+lbl.Comment+"\n"
		}
		name := ""
		if lbl.Name != "" {
			name = lbl.Name+":\n"
		}
		prefix = "\n"+comment+name
	}

	if t.Instruction == nil {
		return fmt.Sprintf("%s[%04X] %02X %-5s : %d%s",
			prefix,
			t.Offset,
			t.Raw,
			"",
			int(t.Raw),
			suffix,
		)
	}

	if len(t.Inline) == 0 {
		return fmt.Sprintf("%s[%04X] %02X %-5s : %s%s",
			prefix,
			t.Offset,
			t.Raw,
			"",
			t.Instruction.String(),
			suffix,
		)
	}

	argstr := []string{}
	for _, a := range t.Inline {
		if lbl, ok := labels[a.Int()]; ok {
			argstr = append(argstr, lbl.Name)
		} else {
			argstr = append(argstr, a.HexString())
		}
	}

	bytestr := []string{}
	for _, a := range t.Inline {
		for _, b := range a.Bytes() {
			//if lbl, ok := labels[a.Int()]; ok {
			//	bytestr = append(bytestr, lbl)
			//} else {
				bytestr = append(bytestr, fmt.Sprintf("%02X", b))
			//}
		}
	}

	switch t.Raw {
	case 0xBB: // push_data
		bs := []byte{}
		for _, val := range t.Inline {
			bs = append(bs, val.Bytes()...)
		}

		return fmt.Sprintf("%s[%04X] %02X (...) : %s %q%s",
			prefix,
			t.Offset,
			t.Raw,
			t.Instruction.String(),
			string(bs[1:len(bs)-1]),
			//strings.Join(argstr[1:], " "),
			suffix,
		)

	//case 0x84, 0x85, 0xBF, 0xC0, // jmp/call


	case 0xC1, 0xEE: // switches
		return fmt.Sprintf("%s[%04X] %02X %-5s : %s %s%s",
			prefix,
			t.Offset,
			t.Raw,
			"",
			t.Instruction.String(),
			strings.Join(argstr, " "),
			suffix,
		)

	default:
		return fmt.Sprintf("%s[%04X] %02X %-5s : %s %s%s",
			prefix,
			t.Offset,
			t.Raw,
			strings.Join(bytestr, " "),
			t.Instruction.String(),
			strings.Join(argstr, " "),
			suffix,
		)

	}

	return fmt.Sprintf("%s%04X: %s %s%s",
		prefix,
		t.Offset,
		t.Instruction.String(),
		strings.Join(argstr, " "),
		suffix,
	)
}

