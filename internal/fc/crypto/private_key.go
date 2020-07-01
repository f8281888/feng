package crypto

//PrivateKey ..
type PrivateKey struct {
	key string
}

//NewPrivateKey ..
func NewPrivateKey(key string) PrivateKey {
	p := PrivateKey{}
	p.key = key
	return p
}
