package script

import (
	"fmt"
)

type Script struct {
	Tokens []*Token
	Warnings []string

	StartAddress int
	StackAddress int

	Labels map[int]*Label
}

type InstrStat struct {
	Instr *Instruction
	Count int
}

func (is InstrStat) String() string {
	return fmt.Sprintf("0x%02X %3d %s", is.Instr.Opcode, is.Count, is.Instr.String())
}

func (s *Script) Stats() Stats {
	st := make(Stats)

	for _, t := range s.Tokens {
		if t.Instruction == nil {
			continue
		}

		op := t.Instruction.Opcode
		if _, ok := st[op]; !ok {
			st[op] = &InstrStat{
				Instr: t.Instruction,
				Count: 0,
			}
		}
		st[op].Count++
	}

	return st
}
