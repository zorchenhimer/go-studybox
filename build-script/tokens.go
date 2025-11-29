package build

import (
	"strings"
	"strconv"
	"fmt"
)

type Token interface {
	Type() string
	String() string
	ValidAfter(t string) bool
	Text() string
}

type TokenDelay struct {
	Value int
	Reset bool
}

func (itm *TokenDelay) Type() string { return "delay" }
func (itm *TokenDelay) String() string { return fmt.Sprintf("{TokenDelay Value:%d Reset:%t}", itm.Value, itm.Reset) }

func (itm *TokenDelay) ValidAfter(t string) bool {
	switch t {
	case "audiooffsets", "pattern", "delay", "page", "script", "tiles":
		return true
	}
	return false
}

func (itm *TokenDelay) Text() string {
	return fmt.Sprintf("delay %d reset:%t", itm.Value, itm.Reset)
}

func parseDelay(key, line string) (Token, error) {
	vals := strings.Split(line, " ")
	itm := &TokenDelay{}

	for _, val := range vals {
		if strings.Contains(val, ":") {
			keyval := strings.SplitN(val, ":", 2)
			if len(keyval) != 2 {
				return nil, fmt.Errorf("Invalid key/value for delay: %s", val)
			}

			if keyval[0] != "reset" {
				return nil, fmt.Errorf("Invalid key/value for delay: %s", val)
			}

			switch strings.ToLower(keyval[1]) {
			case "true", "yes", "1":
				itm.Reset = true
			case "false", "no", "0":
				itm.Reset = false
			default:
				return nil, fmt.Errorf("Invalid reset value: %s", keyval[1])
			}
		} else {
			num, err := strconv.ParseUint(val, 0, 32)
			if err != nil {
				return nil, fmt.Errorf("Invalid delay vaule: %s: %w", val, err)
			}

			itm.Value = int(num)
		}
	}

	return itm, nil
}

type TokenNumValue struct {
	ValType string // delay, padding, page
	Value int
}

func (itm *TokenNumValue) Type() string { return itm.ValType }
func (itm *TokenNumValue) String() string { return fmt.Sprintf("{TokenNumValue ValType:%s Value:%d}", itm.ValType, itm.Value) }

func (itm *TokenNumValue) ValidAfter(t string) bool {
	switch itm.ValType {
	case "page":
		if t == "" {
			return false
		}
		return true
	
	case "padding":
		switch t {
		case "delay", "pattern", "tiles", "script", "page":
			return true
		default:
			return false
		}

	case "version":
		switch t {
		case "", "rom", "fullaudio":
			return true
		default:
			return false
		}
	}

	return false
}

func (itm *TokenNumValue) Text() string {
	return fmt.Sprintf("%s %d", itm.ValType, itm.Value)
}

func parseNumValue(key, val string) (Token, error) {
	v, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("Invalid %s value: %q", key, val)
	}

	return &TokenNumValue{
		ValType: key,
		Value: int(v),
	}, nil
}

type TokenStrValue struct {
	ValType string
	Value string
}

func (itm *TokenStrValue) Type() string { return itm.ValType }
func (itm *TokenStrValue) String() string { return fmt.Sprintf("{TokenStrVal ValType:%s Value:%q}", itm.ValType, itm.Value) }

func (itm *TokenStrValue) ValidAfter(t string) bool {
	switch t {
	case "", "rom", "fullaudio", "version":
		return true
	}
	return false
}

func (itm *TokenStrValue) Text() string {
	return itm.ValType+" "+itm.Value
}

func parseStrValue(key, value string) (Token, error) {
	return &TokenStrValue{
		ValType: key,
		Value: value,
	}, nil
}

type TokenAudioOffsets struct {
	LeadIn uint64
	Data uint64
}

func (itm *TokenAudioOffsets) Type() string { return "audiooffsets" }
func (itm *TokenAudioOffsets) String() string { return fmt.Sprintf("{ItemAudioOffsets LeadIn:%d Data:%d}", itm.LeadIn, itm.Data) }

func (itm *TokenAudioOffsets) ValidAfter(t string) bool {
	switch t {
	case "page", "delay", "script", "tiles", "pattern":
		return true
	}
	return false
}

func (itm *TokenAudioOffsets) Text() string {
	return fmt.Sprintf("audiooffsets leadin:%d data:%d", itm.LeadIn, itm.Data)
}

func parseAudioOffsets(key, line string) (Token, error) {
	vals := strings.Split(line, " ")
	itm := &TokenAudioOffsets{}
	for _, keyval := range vals {
		pair := strings.Split(keyval, ":")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid syntax: %q", keyval)
		}

		num, err := strconv.ParseUint(pair[1], 0, 64)
		if err != nil {
			return nil, err
		}

		switch pair[0] {
		case "leadin":
			itm.LeadIn = num
		case "data":
			itm.Data = num
		default:
			return nil, fmt.Errorf("unknown key: %s", pair[0])
		}
	}

	return itm, nil
}

type TokenData struct {
	ValType string
	Bank int
	Addr int
	File string
}

func (itm *TokenData) Type() string { return itm.ValType }
func (itm *TokenData) String() string {
	return fmt.Sprintf("{TokenData ValType:%s Bank:0x%02X Addr:0x%02X File:%q}",
		itm.ValType,
		itm.Bank,
		itm.Addr,
		itm.File,
	)
}

func (itm *TokenData) ValidAfter(t string) bool {
	switch t {
	case "page", "delay", "script", "tiles", "pattern":
		return true
	}
	return false
}

func (itm *TokenData) Text() string {
	return fmt.Sprintf("%s bank:0x%02X addr:0x%02X file:%q",
		itm.ValType,
		itm.Bank,
		itm.Addr,
		itm.File,
	)
}

func parseData(tokType, vals string) (Token, error) {
	args, err := pKeyVals(vals)
	if err != nil {
		return nil, err
	}

	itm := &TokenData{ValType: tokType}
	for key, value := range args {
		switch key {
		case "bank":
			val, err := strconv.ParseUint(value, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("%s bank value error: %w", key, err)
			}
			itm.Bank = int(val)

		case "addr":
			val, err := strconv.ParseUint(value, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("%s addr value error: %w", key, err)
			}
			itm.Addr = int(val)

		case "file":
			itm.File = value

		default:
			return nil, fmt.Errorf("%s unknown key: %q", tokType, key)
		}
	}

	return itm, nil
}
