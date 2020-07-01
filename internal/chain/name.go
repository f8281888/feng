package chain

import (
	"feng/internal/log"
	"strconv"
)

//Name ..
type Name struct {
	value uint64
}

//ToString ..
func (n *Name) ToString() string {
	b := strconv.FormatUint(n.value, 10)
	return b
}

//Name ..
func (n *Name) Name(input string) {
	b, err := strconv.Atoi(input)
	if err != nil {
		log.AppLog().Errorf("string is error :%s", input)
	}

	n.value = uint64(b)
}
