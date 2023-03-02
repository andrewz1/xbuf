package xbuf

import (
	"encoding/binary"
	"sync"
)

type WB struct {
	b []byte
}

var (
	wbp = sync.Pool{
		New: func() interface{} {
			return &WB{b: make([]byte, 0)}
		},
	}
)

func GetWB() *WB {
	return wbp.Get().(*WB)
}

func PutWB(wb *WB) {
	if wb == nil {
		return
	}
	wb.Reset()
	wbp.Put(wb)
}

func (wb *WB) Reset() {
	wb.b = wb.b[:0]
}

func (wb *WB) PutU8(v byte) {
	wb.b = append(wb.b, v)
}

func (wb *WB) PutU16(v uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	wb.b = append(wb.b, buf...)
}

func (wb *WB) PutU32(v uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	wb.b = append(wb.b, buf...)
}

func (wb *WB) PutU64(v uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	wb.b = append(wb.b, buf...)
}

func (wb *WB) PutBytes(v []byte) {
	wb.b = append(wb.b, v...)
}

func (wb *WB) PutZeros(c int) {
	wb.b = append(wb.b, make([]byte, c)...)
}

func (wb *WB) Len() int {
	return len(wb.b)
}

func (wb *WB) Bytes() []byte {
	return wb.b
}

func (wb *WB) GetBytes() []byte {
	return append([]byte(nil), wb.b...)
}
