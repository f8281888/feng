package chain

import (
	"strings"
)

//Name ..
type Name struct {
	value uint64
}

func charToSymbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64(c-'a') + 6
	}
	if c >= '1' && c <= '5' {
		return uint64(c-'1') + 1
	}
	return 0
}

func stringToUint64(str string) uint64 {
	var n uint64 = 0
	var i int = 0
	max := 12
	if len(str) < 12 {
		max = len(str) - 1
	}

	for ; str[i] != byte(' ') && i < max; i++ {
		n |= charToSymbol(str[i]&0x1f) << (64 - 5*(i+1))
	}

	if i == max {
		n |= charToSymbol(str[max]) & 0xf
	}

	return n
}

//ToString ..
func (n *Name) ToString() string {
	charmap := []byte(".12345abcdefghijklmnopqrstuvwxyz")
	str := []byte(".............")
	var tmp uint64 = n.value
	for i := 0; i <= 12; i++ {
		var tmpByte byte
		var index uint32
		if i == 0 {
			tmpByte = 0xf
			index = 4
		} else {
			tmpByte = 0x1f
			index = 5
		}

		var c byte = charmap[tmp&uint64(tmpByte)]
		str[12-i] = c
		tmp >>= index
	}

	s := string(str)
	strings.TrimRight(s, ".")
	return s
}

//Name ..
func (n *Name) Name(input string) {
	n.value = stringToUint64(input)
}
