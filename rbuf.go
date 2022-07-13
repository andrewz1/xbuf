package xbuf

import (
	"encoding/binary"
	"sync"
)

type RB struct {
	b []byte // buf data
	p int    // buff ptr
	l int    // bytes left
	n bool   // nested buffer flag
}

var (
	rbp = sync.Pool{
		New: func() interface{} {
			return &RB{b: make([]byte, 0)}
		},
	}
)

func GetRB(b []byte) *RB {
	rb := rbp.Get().(*RB)
	if len(b) > 0 {
		rb.b = append(rb.b, b...)
		rb.l = len(b)
	}
	return rb
}

func PutRB(rb *RB) {
	if rb == nil || rb.n {
		return
	}
	rb.Reset()
	rbp.Put(rb)
}

func (rb *RB) Reset() {
	rb.b = rb.b[:0]
	rb.p = 0
	rb.l = 0
}

func (rb *RB) Len() int {
	return len(rb.b)
}

func (rb *RB) Ptr() int {
	return rb.p
}

func (rb *RB) Left() int {
	return rb.l
}

func (rb *RB) Append(b []byte) {
	if len(b) == 0 {
		return
	}
	rb.b = append(rb.b, b...)
	rb.l += len(b)
}

func (rb *RB) Set(b []byte) {
	rb.b = append(rb.b[:0], b...)
	rb.p = 0
	rb.l = len(b)
}

func (rb *RB) shift(n int) {
	rb.p += n
	rb.l -= n
}

func (rb *RB) save() (p, l int) {
	p = rb.p
	l = rb.l
	return
}

func (rb *RB) restore(p, l int) {
	rb.p = p
	rb.l = l
}

func (rb *RB) GetU8() (byte, bool) {
	if rb.l < 1 {
		return 0, false
	}
	defer rb.shift(1)
	return rb.b[rb.p], true
}

func (rb *RB) MustGetU8() byte {
	rv, _ := rb.GetU8()
	return rv
}

func (rb *RB) GetU16() (uint16, bool) {
	if rb.l < 2 {
		return 0, false
	}
	defer rb.shift(2)
	return binary.BigEndian.Uint16(rb.b[rb.p:]), true
}

func (rb *RB) MustGetU16() uint16 {
	rv, _ := rb.GetU16()
	return rv
}

func (rb *RB) GetU24() (uint32, bool) {
	if rb.l < 3 {
		return 0, false
	}
	defer rb.shift(3)
	var t [4]byte
	copy(t[1:], rb.b[rb.p:])
	return binary.BigEndian.Uint32(t[:]), true
}

func (rb *RB) MustGetU24() uint32 {
	rv, _ := rb.GetU24()
	return rv
}

func (rb *RB) GetU32() (uint32, bool) {
	if rb.l < 4 {
		return 0, false
	}
	defer rb.shift(4)
	return binary.BigEndian.Uint32(rb.b[rb.p:]), true
}

func (rb *RB) MustGetU32() uint32 {
	rv, _ := rb.GetU32()
	return rv
}

func (rb *RB) GetU64() (uint64, bool) {
	if rb.l < 8 {
		return 0, false
	}
	defer rb.shift(8)
	return binary.BigEndian.Uint64(rb.b[rb.p:]), true
}

func (rb *RB) MustGetU64() uint64 {
	rv, _ := rb.GetU64()
	return rv
}

func (rb *RB) bytes(n int) []byte {
	switch {
	case n < 0:
		return nil
	case n == 0:
		return make([]byte, 0)
	default:
		defer rb.shift(n)
		return rb.b[rb.p : rb.p+n]
	}
}

func (rb *RB) Bytes(n int) []byte {
	if rb.l < n {
		return nil
	}
	return rb.bytes(n)
}

func (rb *RB) GetBytes(n int) ([]byte, bool) {
	if rb.l < n {
		return nil, false
	}
	return append([]byte{}, rb.bytes(n)...), true
}

func (rb *RB) Skip(n int) bool {
	if rb.l < n {
		return false
	}
	rb.shift(n)
	return true
}

func (rb *RB) SkipL8() (ok bool) {
	p, l := rb.save()
	defer func() {
		if !ok {
			rb.restore(p, l)
		}
	}()
	var n byte
	if n, ok = rb.GetU8(); !ok {
		return
	}
	ok = rb.Skip(int(n))
	return
}

func (rb *RB) SkipL16() (ok bool) {
	p, l := rb.save()
	defer func() {
		if !ok {
			rb.restore(p, l)
		}
	}()
	var n uint16
	if n, ok = rb.GetU16(); !ok {
		return
	}
	ok = rb.Skip(int(n))
	return
}

func (rb *RB) nestedRB(n int) *RB {
	return &RB{
		b: rb.bytes(n),
		l: n,
		n: true,
	}
}

func (rb *RB) GetNested(n int) (nb *RB, ok bool) {
	if rb.l < n {
		return
	}
	nb = rb.nestedRB(n)
	ok = true
	return
}

func (rb *RB) GetNestedL8() (nb *RB, ok bool) {
	p, l := rb.save()
	defer func() {
		if !ok {
			rb.restore(p, l)
			nb = nil
		}
	}()
	var n byte
	if n, ok = rb.GetU8(); !ok {
		return
	}
	nb, ok = rb.GetNested(int(n))
	return
}

func (rb *RB) GetNestedL16() (nb *RB, ok bool) {
	p, l := rb.save()
	defer func() {
		if !ok {
			rb.restore(p, l)
			nb = nil
		}
	}()
	var n uint16
	if n, ok = rb.GetU16(); !ok {
		return
	}
	nb, ok = rb.GetNested(int(n))
	return
}

func (rb *RB) GetBuf(n int) []byte {
	if n > len(rb.b) {
		return nil
	}
	return append([]byte{}, rb.b[:n]...)
}

func (rb *RB) String() string {
	return string(rb.b)
}
