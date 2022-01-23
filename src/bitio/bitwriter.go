package bitio

import (
	"fmt"
	"io/ioutil"
)

type BitWriter struct {
	stream  []byte
	buffer  byte
	bufSize int
}

func NewBitWriter() *BitWriter {
	return &BitWriter{
		[]byte{},
		0x00,
		0,
	}
}

func (bw *BitWriter) WriteToFile(fileName string) {
	// Finalize buffer
	if bw.bufSize > 0 {
		bw.stream = append(bw.stream, bw.buffer)
		bw.buffer = 0x00
		bw.bufSize = 0
	}

	err := ioutil.WriteFile(fileName, bw.stream, 0644)

	if err != nil {
		fmt.Printf("File write error: %s\n", err)
	}
}

func (bw *BitWriter) WriteStr(str string) {
	bw.stream = append(bw.stream, str...)
}

func (bw *BitWriter) WriteBits(bits int, num int) {
	for i := num - 1; i >= 0; i-- {
		bw.WriteBit(!((1<<i)&bits == 0))
	}
}

func (bw *BitWriter) WriteBit(bit bool) {
	if bw.bufSize == 8 {
		bw.stream = append(bw.stream, bw.buffer)
		bw.buffer = 0x00
		bw.bufSize = 0
	}

	if bit {
		bw.buffer |= 1 << (8 - bw.bufSize - 1)
	}

	bw.bufSize++
}
