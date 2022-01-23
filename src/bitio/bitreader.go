package bitio

import (
	"fmt"
	"io/ioutil"
)

type BitReader struct {
	stream  string
	buffer  byte
	bufSize int
}

func NewBitReader(fileName string) *BitReader {
	fileBuf, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("File read error: %s\n", err)
		return nil
	}

	return &BitReader{
		string(fileBuf),
		0x00,
		0,
	}
}

func (br *BitReader) HasNext() bool {
	return len(br.stream) > 0
}

func (br *BitReader) ReadBit() bool {
	if br.bufSize == 0 {
		br.buffer = br.ReadByte()
		br.bufSize = 8
	}

	bit := !(br.buffer&(1<<(br.bufSize-1)) == 0)

	br.bufSize--

	return bit
}

func (br *BitReader) ReadBits(numBits int) int {
	bits := 0

	for i := 0; i < numBits; i++ {
		bits <<= 1
		if br.ReadBit() {
			bits |= 1
		}
	}

	return bits
}

func (br *BitReader) PeekByte() byte {
	if len(br.stream) > 0 {
		return br.stream[0]
	}

	return 0x00
}

func (br *BitReader) ReadByte() byte {
	if len(br.stream) == 0 {
		return 0x00
	}

	byteRead := br.stream[0]
	br.stream = br.stream[1:]

	return byteRead
}
