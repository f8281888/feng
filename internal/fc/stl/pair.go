package stl

//Pair ..
type Pair struct {
	First  Item
	Second Item
}

//SetFirst ..
func (p *Pair) SetFirst(i Item) {
	p.First = i
}

//SetSecond ..
func (p *Pair) SetSecond(i Item) {
	p.Second = i
}
