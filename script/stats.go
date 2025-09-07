package script

import (
	"io"
	"slices"
	"maps"
	"fmt"
)

type Stats map[byte]*InstrStat

func (this Stats) Add(that Stats) {
	for _, st := range that {
		op := st.Instr.Opcode
		if _, ok := this[op]; !ok {
			this[op] = st
		} else {
			this[op].Count += that[op].Count
		}
	}
}

func (s Stats) WriteTo(w io.Writer) (int64, error) {
	count := int64(0)
	keys := slices.Sorted(maps.Keys(s))

	unknownInstr := 0
	unknownUses := 0

	for _, key := range keys {
		n, err := fmt.Fprintln(w, s[key])
		count += int64(n)
		if err != nil {
			return count, err
		}

		if s[key].Instr.Name == "" {
			unknownInstr++
			unknownUses += s[key].Count
		}
	}

	n, err := fmt.Fprintln(w, "\nUnknown uses:", unknownUses)
	count += int64(n)
	if err != nil {
		return count, err
	}

	n, err = fmt.Fprintln(w, "Unknown instructions:", unknownInstr)
	count += int64(n)
	if err != nil {
		return count, err
	}

	return count, nil
}
