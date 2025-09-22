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

	cdl string // CDL string type

	Instruction *Instruction
}

func (t Token) String(labels map[int]*Label, suppAddr bool) string {
	suffix := ""
	switch t.Raw {
	case 0x86, 0xAC, 0xAA: // Newline after return, long_return, & long_jump
		suffix = "\n"
	}

	offset := ""
	if !suppAddr {
		offset = fmt.Sprintf("[%04X] ", t.Offset)
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

	if t.IsVariable {
		return fmt.Sprintf("%s%s%02X %-5s : %d %s%s",
			prefix,
			offset,
			t.Raw,
			"",
			t.Inline[0].Int(),
			t.Inline[0].HexString(),
			suffix,
		)
	}

	if t.Instruction == nil {
		if t.IsData == false {
			return fmt.Sprintf("%s%s%02X %-5s : %d%s",
				prefix,
				offset,
				t.Raw,
				"",
				int(t.Raw),
				suffix,
			)
		} else if t.IsData {
			return fmt.Sprintf("%s%s%02X %-5s : %d%s",
				prefix,
				offset,
				t.Raw,
				t.cdl,
				int(t.Raw),
				suffix,
			)
		}
	}

	if len(t.Inline) == 0 {
		return fmt.Sprintf("%s%s%02X %-5s : %s%s",
			prefix,
			offset,
			t.Raw,
			"",
			t.Instruction.String(),
			suffix,
		)
	}

	argstr := []string{}
	for _, a := range t.Inline {
		if lbl, ok := labels[a.Int()]; ok && !t.Instruction.InlineImmediate {
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
		raw := []byte{}

		ascii := true
		for _, val := range t.Inline[1:len(t.Inline)-1] {
			for _, b := range val.Bytes() {
				raw = append(raw, b)
				if b < 0x20 || b > 0x7E {
					ascii = false
				}
			}
		}

		bs := ""
		if ascii {
			bs = fmt.Sprintf("%q", string(raw))
		} else {
			vals := []string{}
			for _, b := range raw {
				if b >= 0x20 && b <= 0x7E {
					vals = append(vals, fmt.Sprintf("0x%02X{%c}", b, b))
				} else {
					vals = append(vals, fmt.Sprintf("0x%02X", b))
				}
			}
			bs = "["+strings.Join(vals, " ")+"]"
		}

		//for _, val := range t.Inline {
		//	//bs = append(bs, val.Bytes()...)
		//	for _, b := range val.Bytes() {
		//		// These strings are strictly binary or ascii.  If there's
		//		// non-ascii, don't try and read it as unicode if we find
		//		// some "valid" code points.  Eg, 0xD?, 0xB? (%110?_????, %10??_????)
		//		if b < 0x20 || b > 0x7E {
		//			bs = append(bs, fmt.Sprintf("\\x%02x", b))
		//		} else {
		//			bs = append(bs, string(b))
		//		}
		//	}
		//}

		return fmt.Sprintf("%s%s%02X (...) : %s %s%s",
			prefix,
			offset,
			t.Raw,
			t.Instruction.String(),
			//string(bs[1:len(bs)-1]),
			//strings.Join(bs[1:len(bs)-1], ""),
			//strings.Join(argstr[1:], " "),
			bs,
			suffix,
		)

	//case 0x84, 0x85, 0xBF, 0xC0, // jmp/call


	case 0xC1, 0xEE: // switches
		return fmt.Sprintf("%s%s%02X %-5s : %s %s%s",
			prefix,
			offset,
			t.Raw,
			"",
			t.Instruction.String(),
			strings.Join(argstr, " "),
			suffix,
		)

	default:
		return fmt.Sprintf("%s%s%02X %-5s : %s %s%s",
			prefix,
			offset,
			t.Raw,
			strings.Join(bytestr, " "),
			t.Instruction.String(),
			strings.Join(argstr, " "),
			suffix,
		)

	}

	return fmt.Sprintf("%s%s%s %s%s",
		prefix,
		offset,
		t.Instruction.String(),
		strings.Join(argstr, " "),
		suffix,
	)
}

