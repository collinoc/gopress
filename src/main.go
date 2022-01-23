package main

// Compressed files are output to the same directory as <filename_and_ext>.lzw
// Decompressed files are output to the same directory and overwrites any file with the same name

import (
	"errors"
	"fmt"
	"math"
	"os"
	"gopress/bitio"
)

type Mode int

const (
	Compress Mode = iota
	Decompress
)

const startingCodeLength = 9

func main() {
	mode, file, err := handleCommandLine()

	if err != nil {
		fmt.Println(os.Args)
		fmt.Println(err)
		return
	}

	switch mode {
	case Compress:
		compressFile(file)
		break
	case Decompress:
		decompressFile(file)
		break
	}
}

func toStr(b byte) string {
	return string([]byte{b})
}

func compressFile(file string) {
	fmt.Print("Compressing " + file + "... ")

	br := bitio.NewBitReader(file)
	bw := bitio.NewBitWriter()
	codeBitLen := startingCodeLength
	dict := possibleBytes()

	if br == nil {
		return
	}

	for br.HasNext() {
		if len(dict) == int(math.Pow(2, float64(codeBitLen))) {
			codeBitLen++

			if codeBitLen > 12 {
				codeBitLen = 9
				dict = possibleBytes()
			}
		}

		currStr := toStr(br.ReadByte())

		// Todo: Increase performance of contains() check
		for br.HasNext() && contains(dict, currStr+toStr(br.PeekByte())) {
			currStr += toStr(br.ReadByte())
		}

		dict = append(dict, currStr+toStr(br.PeekByte()))
		bw.WriteBits(indexOf(dict, currStr), codeBitLen)
	}

	bw.WriteToFile(file + ".lzw")
	fmt.Println("Done.")
}

func decompressFile(file string) {
	fmt.Print("Decompressing " + file + "... ")

	br := bitio.NewBitReader(file)
	bw := bitio.NewBitWriter()
	codeBitLen := startingCodeLength
	previousWritten := ""
	dict := possibleBytes()

	if br == nil {
		return
	}

	reset := false

	for br.HasNext() {
		if reset {
			dict = possibleBytes()
			reset = false
		}

		if len(dict) == int(math.Pow(2, float64(codeBitLen)))-1 {
			codeBitLen++

			if codeBitLen > 12 {
				codeBitLen = 9
				reset = true
			}
		}

		symbol := br.ReadBits(codeBitLen)

		if symbol < len(dict) {
			curr := dict[symbol]
			bw.WriteStr(curr)

			if previousWritten != "" {
				dict = append(dict, previousWritten+toStr(curr[0]))
			}

			previousWritten = curr
		} else {
			curr := previousWritten + toStr(previousWritten[0])
			dict = append(dict, curr)
			bw.WriteStr(curr)
			previousWritten = curr
		}
	}

	// Trim off .lzw
	bw.WriteToFile(file[:len(file)-3])

	fmt.Println("Done.")
}

func indexOf(a []string, item string) int {
	for idx, str := range a {
		if item == str {
			return idx
		}
	}
	return -1
}

func contains(a []string, item string) bool {
	for _, str := range a {
		if item == str {
			return true
		}
	}
	return false
}

func possibleBytes() []string {
	bytes := []string{}

	for n := 0; n < 256; n++ {
		bytes = append(bytes, toStr(byte(n)))
	}
	return bytes
}

func handleCommandLine() (Mode, string, error) {
	if len(os.Args) != 3 {
		return 0, "", errors.New("Bad command line argument length. Use format <-c | -d> \"<filename>\"")
	}

	mode := Compress

	if os.Args[1] == "-d" {
		mode = Decompress
	} else if os.Args[1] != "-c" {
		return 0, "", errors.New("Bad mode format. Use format <-c | -d> \"<filename>\"")
	}

	return mode, os.Args[2], nil
}

// type BitWriter struct {
// 	stream  []byte
// 	buffer  byte
// 	bufSize int
// }

// func NewBitWriter() *BitWriter {
// 	return &BitWriter{
// 		[]byte{},
// 		0x00,
// 		0,
// 	}
// }

// func (bw *BitWriter) writeToFile(fileName string) {
// 	// Finalize buffer
// 	if bw.bufSize > 0 {
// 		bw.stream = append(bw.stream, bw.buffer)
// 		bw.buffer = 0x00
// 		bw.bufSize = 0
// 	}

// 	err := ioutil.WriteFile(fileName, bw.stream, 0644)

// 	if err != nil {
// 		fmt.Printf("File write error: %s\n", err)
// 	}
// }

// func (bw *BitWriter) writeStr(str string) {
// 	bw.stream = append(bw.stream, str...)
// }

// func (bw *BitWriter) writeBits(bits int, num int) {
// 	for i := num - 1; i >= 0; i-- {
// 		bw.writeBit(!((1<<i)&bits == 0))
// 	}
// }

// func (bw *BitWriter) writeBit(bit bool) {
// 	if bw.bufSize == 8 {
// 		bw.stream = append(bw.stream, bw.buffer)
// 		bw.buffer = 0x00
// 		bw.bufSize = 0
// 	}

// 	if bit {
// 		bw.buffer |= 1 << (8 - bw.bufSize - 1)
// 	}

// 	bw.bufSize++
// }

// type BitReader struct {
// 	stream  string
// 	buffer  byte
// 	bufSize int
// }

// func NewBitReader(fileName string) *BitReader {
// 	fileBuf, err := ioutil.ReadFile(fileName)

// 	if err != nil {
// 		fmt.Printf("File read error: %s\n", err)
// 		return nil
// 	}

// 	return &BitReader{
// 		string(fileBuf),
// 		0x00,
// 		0,
// 	}
// }

// func (br *BitReader) hasNext() bool {
// 	return len(br.stream) > 0
// }

// func (br *BitReader) readBit() bool {
// 	if br.bufSize == 0 {
// 		br.buffer = br.readByte()
// 		br.bufSize = 8
// 	}

// 	bit := !(br.buffer&(1<<(br.bufSize-1)) == 0)

// 	br.bufSize--

// 	return bit
// }

// func (br *BitReader) readBits(numBits int) int {
// 	bits := 0

// 	for i := 0; i < numBits; i++ {
// 		bits <<= 1
// 		if br.readBit() {
// 			bits |= 1
// 		}
// 	}

// 	return bits
// }

// func (br *BitReader) peekByte() byte {
// 	if len(br.stream) > 0 {
// 		return br.stream[0]
// 	}

// 	return 0x00
// }

// func (br *BitReader) readByte() byte {
// 	if len(br.stream) == 0 {
// 		return 0x00
// 	}

// 	byteRead := br.stream[0]
// 	br.stream = br.stream[1:]

// 	return byteRead
// }
