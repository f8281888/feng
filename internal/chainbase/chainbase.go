package chainbase

//ChainBase ..
type ChainBase struct {
	rversion        uint64
	sizeOfValueType uint32
	sizeOfThis      uint32
}

//GetRversion ..
func (c ChainBase) GetRversion() uint64 {
	return c.rversion
}

//SetRversion ..
func (c *ChainBase) SetRversion(u uint64) {
	c.rversion = u
}

//Undo ..
func (c ChainBase) Undo() {
	println("ChainBase Undo")
}
