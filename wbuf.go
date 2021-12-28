package xbuf

import (
	"encoding/binary"

	"github.com/valyala/bytebufferpool"
)

type WB struct {
	*bytebufferpool.ByteBuffer
}

var (
	wbp bytebufferpool.Pool
)

func GetWB() *WB {
	return &WB{
		ByteBuffer: wbp.Get(),
	}
}

func PutWB(wb *WB) {
	if wb == nil {
		return
	}
	wbp.Put(wb.ByteBuffer)
}

func (wb *WB) PutU8(v byte) {
	wb.B = append(wb.B, v)
}

func (wb *WB) PutU16(v uint16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	wb.B = append(wb.B, buf[:]...)
}

func (wb *WB) PutU32(v uint32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	wb.B = append(wb.B, buf[:]...)
}

func (wb *WB) PutU64(v uint64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	wb.B = append(wb.B, buf[:]...)
}

func (wb *WB) PutBytes(v []byte) {
	wb.B = append(wb.B, v...)
}

func (wb *WB) PutZeros(c int) {
	wb.B = append(wb.B, make([]byte, c)...)
}
