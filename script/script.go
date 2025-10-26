package script

import (
	"fmt"
	"os"
)

type Script struct {
	Tokens []*Token
	Warnings []string

	StartAddress int
	StackAddress int

	Labels map[int]*Label
	CDL *CodeDataLog

	origSize int // size of the binary input
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

func (s *Script) DebugCDL(filename string) error {
	if s.origSize == 0 {
		return fmt.Errorf("origSize == 0")
	}

	//if s.CDL.cache == nil {
	//	err := s.CDL.doCache()
	//	if err != nil {
	//		return fmt.Errorf("doCache() error: %w", err)
	//	}
	//}

	dat := make([]byte, s.origSize)
	for i := 2; i < len(dat); i++ {
		if val, ok := s.CDL.cache[i+0x6000]; ok {
			dat[i] = byte(val)
		}
	}

	err := os.WriteFile(filename, dat, 0644)
	if err != nil {
		return fmt.Errorf("WriteFile() error: %w", err)
	}

	return nil
}
