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
	for idx := range a {
		if item == a[idx] {
			return idx
		}
	}
	return -1
}

func contains(a []string, item string) bool {
	for idx := range a {
		if item == a[idx] {
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