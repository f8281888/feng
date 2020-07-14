package crypto

//PublicKey ..
type PublicKey struct {
	key string
}

//NewPublicKey ..
func NewPublicKey(key string) PublicKey {
	p := PublicKey{}
	p.key = key
	return p
}

//ToString ..
func (p PublicKey) ToString() string {
	return p.key
}
