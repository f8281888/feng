package common

import (
	"bytes"
	"encoding/binary"
)

func ternary(expr bool, whenTrue, whenFalse interface{}) interface{} {
	if expr == true {
		return whenTrue
	}

	return whenFalse
}

//Min 小
func Min(a, b int) int {
	i := ternary(a <= b, a, b)
	val, _ := i.(int)
	return val
}

//Max 大
func Max(a, b int) int {
	i := ternary(a >= b, a, b)
	val, _ := i.(int)
	return val
}

//MinUint64 小
func MinUint64(a, b uint64) uint64 {
	i := ternary(a <= b, a, b)
	val, _ := i.(uint64)
	return val
}

//MaxUint64 大
func MaxUint64(a, b uint64) uint64 {
	i := ternary(a >= b, a, b)
	val, _ := i.(uint64)
	return val
}

//BytesToInt ..
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int(x)
}

//IntToBytes ..
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//BoolToBytes ..
func BoolToBytes(b bool) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, b)
	return bytesBuffer.Bytes()
}

//BytesToBool ..
func BytesToBool(b []byte) bool {
	bytesBuffer := bytes.NewBuffer(b)
	var x bool
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

//BytesToUint64 ..
func BytesToUint64(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return uint64(x)
}
