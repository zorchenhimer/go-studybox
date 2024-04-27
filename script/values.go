package script

import (
	"fmt"
)

type InlineVal interface {
	HexString() string
	Bytes() []byte
	Int() int
}

type ByteVal byte

func (bv ByteVal) HexString() string {
	return fmt.Sprintf("$%02X", bv)
}

func (bv ByteVal) Bytes() []byte {
	return []byte{byte(bv)}
}

func (bv ByteVal) Int() int {
	return int(bv)
}

type WordVal [2]byte

func NewWordVal(v []byte) WordVal {
	if len(v) != 2 {
		panic("WordVal must be two bytes")
	}

	return WordVal([2]byte{v[0], v[1]})
}

func (wv WordVal) HexString() string {
	return fmt.Sprintf("$%02X%02X", wv[1], wv[0])
}

func (wv WordVal) Bytes() []byte {
	return []byte{wv[0], wv[1]}
}

func (wv WordVal) Int() int {
	return (int(wv[1]) << 8) | int(wv[0])
}
