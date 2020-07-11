package io

import (
	"feng/internal/fc/common"
	"feng/internal/fc/network"
)

//Steam ..
type Steam interface {
}

//Unpack ..
func Unpack(s *network.MbPeekDatastream, vi uint32) {
	var v uint64 = 0
	var b []byte
	var by uint8 = 0

	for {
		s.Get(&b)
		a := b[0] & 0x7f
		b[0] = a
		v |= uint64(common.BytesToInt(b)) << by
		by += 7

		b[0] &= 0x80
		if common.BytesToBool(b) && by < 32 {
			break
		}
	}

	vi = uint32(v)
}
