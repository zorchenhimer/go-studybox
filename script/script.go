package script

import (
)

type Script struct {
	Tokens []*Token
	Warnings []string

	StartAddress int
	StackAddress int

	Labels map[int]*Label
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
