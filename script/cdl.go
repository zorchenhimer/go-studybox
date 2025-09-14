package script

import (
	"io"
	"os"
	"encoding/json"
	"strconv"
	"fmt"
	"slices"
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
	//cdlOpCode  cdlBit = 0x04
)

func (cdl *CodeDataLog) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, werr := cdl.WriteTo(file)
	err = file.Close()
	if err != nil {
		return err
	}

	return werr
}

func getRanges(list []int) []CdlRange {
	//fmt.Printf("getRanges(%v)\n", list)
	data := []CdlRange{}

	start := -1
	//end := -1
	prev := -1
	for _, addr := range list {
		if start == -1 {
			start = addr
		}

		if prev != -1 && prev != addr-1 {
			data = append(data, CdlRange{
				Start: fmt.Sprintf("0x%X", start),
				End: fmt.Sprintf("0x%X", prev),
			})

			//fmt.Printf("start: 0x%X end: 0x%X\n", start, prev)

			start = addr
		}

		prev = addr
	}

	if start != -1 && prev != -1 {
		data = append(data, CdlRange{
			Start: fmt.Sprintf("0x%X", start),
			End: fmt.Sprintf("0x%X", prev),
		})

		//fmt.Println("start:", start, "end:", prev)
	}

	return data
}

func (cdl *CodeDataLog) WriteTo(w io.Writer) (int64, error) {
	clean := &CodeDataLog{
		Code: []CdlRange{},
		Data: []CdlRange{},
	}

	keys := []int{}
	for k, _ := range cdl.cache {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	code := []int{}
	data := []int{}

	for _, k := range keys {
		b := cdl.cache[k]
		if b & cdlCode == cdlCode {
			code = append(code, k)
		}

		if b & cdlData == cdlData {
			data = append(data, k)
		}
	}

	clean.Code = getRanges(code)
	clean.Data = getRanges(data)

	raw, err := json.MarshalIndent(clean, "", "\t")
	if err != nil {
		return 0, err
	}

	n, err := w.Write(raw)
	return int64(n), err
}

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
