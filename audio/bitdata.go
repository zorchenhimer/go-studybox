package audio

type BitData struct {
	data    []byte
	next    int
	current byte
	left    int // bits left in current
}

func NewBitData(data []byte) *BitData {
	if len(data) == 0 {
		panic("no data")
	}

	return &BitData{
		data: data,
		next: 1,
		current: data[0],
		left: 7,
	}
}

// Returns the bit in the lowest position, and end of data.  false if nothing left.
func (b *BitData) Next() (byte, bool) {
	if b.left < 0 {
		if len(b.data) <= b.next {
			return 0, false
		}

		b.current = b.data[b.next]
		b.next++
		b.left = 7
		return 0, true
	}

	ret := (b.current >> b.left) & 0x01
	b.left--
	return ret, true
}

func (b *BitData) Peek() byte {
	left := b.left
	current := b.current

	if left < 0 {
		if len(b.data) <= b.next {
			return 0
		}
		current = b.data[b.next]
		left = 7
	}

	return (current >> left) & 0x01
}

