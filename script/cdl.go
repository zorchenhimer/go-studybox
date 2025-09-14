package script

import (
	"io"
	"os"
	"encoding/json"
	"strconv"
	"fmt"
)

type CodeDataLog struct {
	Code []CdlRange
	Data []CdlRange

	cache map[int]cdlBit
	offset int
}

type CdlRange struct {
	// strings cuz json doesn't know wtf hexadecimal is
	Start string
	End   string
}

type cdlBit byte

var (
	cdlUnknown cdlBit = 0x00
	cdlCode    cdlBit = 0x01
	cdlData    cdlBit = 0x02
)

func (cdl *CodeDataLog) setData(scriptOffset int) {
	if cdl.cache == nil {
		err := cdl.doCache()
		if err != nil {
			panic(fmt.Sprintf("CDL data error: %w", err))
		}
	}

	cdl.cache[scriptOffset+cdl.offset] |= cdlData
}

func (cdl *CodeDataLog) setCode(scriptOffset int) {
	if cdl.cache == nil {
		err := cdl.doCache()
		if err != nil {
			panic(fmt.Sprintf("CDL data error: %w", err))
		}
	}

	cdl.cache[scriptOffset+cdl.offset] |= cdlCode
}

func (cdl *CodeDataLog) doCache() error {
	cdl.cache = make(map[int]cdlBit)

	for _, rng := range cdl.Code {
		start, err := strconv.ParseInt(rng.Start, 0, 32)
		if err != nil {
			return fmt.Errorf("Invalid start: %q", rng.Start)
		}

		end, err := strconv.ParseInt(rng.End, 0, 32)
		if err != nil {
			return fmt.Errorf("Invalid end: %q", rng.End)
		}

		for i := int(start); i <= int(end); i++ {
			if _, ok := cdl.cache[i]; !ok {
				cdl.cache[i] = cdlUnknown
			}

			cdl.cache[i] |= cdlCode
		}
	}

	for _, rng := range cdl.Data {
		start, err := strconv.ParseInt(rng.Start, 0, 32)
		if err != nil {
			return fmt.Errorf("Invalid start: %q", rng.Start)
		}

		end, err := strconv.ParseInt(rng.End, 0, 32)
		if err != nil {
			return fmt.Errorf("Invalid end: %q", rng.End)
		}

		for i := int(start); i <= int(end); i++ {
			if _, ok := cdl.cache[i]; !ok {
				cdl.cache[i] = cdlUnknown
			}

			cdl.cache[i] |= cdlData
		}
	}

	return nil
}

func CdlFromJson(r io.Reader) (*CodeDataLog, error) {
	cdl := &CodeDataLog{}
	dec := json.NewDecoder(r)
	err := dec.Decode(cdl)
	if err != nil {
		return nil, err
	}

	return cdl, nil
}

func CdlFromJsonFile(filename string) (*CodeDataLog, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return CdlFromJson(file)
}

func (cdl *CodeDataLog) IsData(addr int) bool {
	if cdl.cache == nil {
		err := cdl.doCache()
		if err != nil {
			panic(fmt.Sprintf("CDL data error: %w", err))
		}
	}

	val, ok := cdl.cache[addr]
	if !ok {
		return false
	}

	return val & cdlData == cdlData
}

func (cdl *CodeDataLog) IsCode(addr int) bool {
	if cdl.cache == nil {
		err := cdl.doCache()
		if err != nil {
			panic(fmt.Sprintf("CDL data error: %w", err))
		}
	}

	val, ok := cdl.cache[addr]
	if !ok {
		return false
	}

	return val & cdlCode == cdlCode
}
