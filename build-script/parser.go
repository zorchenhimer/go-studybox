package build

import (
	"fmt"
	"bufio"
	"os"
	"io"
	"testing"
	"strings"
	"unicode"
)

type parseFunc func(key, values string) (Token, error)

var parseFuncs = map[string]parseFunc{
	"rom":          parseStrValue,
	"fullaudio":    parseStrValue,
	"audiooffsets": parseAudioOffsets,

	"page":    parseNumValue,
	"padding": parseNumValue,
	"version": parseNumValue,

	"delay": parseDelay,

	"script":  parseData,
	"tiles":   parseData,
	"pattern": parseData,
}

func ParseFile(filename string) ([]Token, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

func Parse(r io.Reader) ([]Token, error) {
	items := []Token{}
	scanner := bufio.NewScanner(r)
	prev := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// blanks and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		splitln := strings.SplitN(line, " ", 2)
		if len(splitln) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}

		var itm Token
		var err error

		// TODO: Some of these will need to be unique in the file
		// (rom, version, fullaudio), and probably exclusive/ordered.
		// IE, rom, version, & fullaudio must come before the first
		// page ##, and must only appear once.  if a page has an
		// audio line, all pages must have one and there can be no
		// fullaudio.

		if fn, ok := parseFuncs[splitln[0]]; ok {
			itm, err = fn(splitln[0], splitln[1])
		} else {
			return nil, fmt.Errorf("unknown line: %s\n", splitln[0])
		}

		if err != nil {
			return nil, err
		}

		if !itm.ValidAfter(prev) {
			if prev == "" {
				prev = "[empty]"
			}
			return nil, fmt.Errorf("%s not valid after %s", itm.Type(), prev)
		}
		prev = itm.Type()

		items = append(items, itm)
	}

	return items, nil
}

var t *testing.T

func pKeyVals(values string) (map[string]string, error) {
	m := map[string]string{}

	start := 0
	runes := []rune(values)
	currentKey := ""
	quote := false
	tlog("[pKeyVals] start")
	var i int

	for i = 0; i < len(runes); i++ {
		tlogf("[pKeyVals] rune: %c\n", runes[i])
		if !quote && runes[i] == '"' {
			quote = true
			tlog("[pKeyVals] start quote")
			continue
		}

		if !quote && unicode.IsSpace(runes[i]) {
			tlog("[pKeyVals] !quote && IsSpace()")
			if currentKey == "" {
				tlog("[pKeyVals] currentKey empty")
				start = i+1
			} else {
				tlog("[pKeyVals] currentKey not empty")
				m[currentKey] = string(runes[start:i])
				currentKey = ""
				start = i+1
			}
			continue
		}

		if quote && runes[i] == '"' {
			tlog("[pKeyVals] quote && rune == \"")
			quote = false

			if currentKey == "" {
				currentKey = string(runes[start+1:i])
				i++
				start = i+1
				tlogf("[pKeyVals] currentKey empty; set to %s\n", currentKey)

			} else {
				m[currentKey] = string(runes[start+1:i])
				currentKey = ""
				tlog("[pKeyVals] currentKey not empty")
			}
			continue
		}

		if runes[i] == ':' {
			tlog("[pKeyVals] rune == :")
			if i == start {
				return nil, fmt.Errorf("missing key")
			}

			currentKey = string(runes[start:i])
			start = i+1
			continue
		}
	}

	if quote {
		return nil, fmt.Errorf("missmatched quote")
	}

	if currentKey != "" {
		//return nil, fmt.Errorf("missing value for %q", currentKey)
		tlogf("[pKeyVals] outside loop assign m[%s] = %s\n", currentKey, string(runes[start:i]))
		m[currentKey] = string(runes[start:i])
	}

	return m, nil
}

func tlog(args ...any) {
	if t != nil {
		t.Log(args...)
	}
}

func tlogf(format string, args ...any) {
	if t != nil {
		t.Logf(format, args...)
	}
}
