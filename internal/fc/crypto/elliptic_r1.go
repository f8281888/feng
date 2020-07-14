package crypto

//PrivateKey ..
type PrivateKey struct {
	key string
}

//Tostring ..
func (p *PrivateKey) Tostring() string {
	return p.key
}

//NewPrivateKey ..
func NewPrivateKey(key string) PrivateKey {
	p := PrivateKey{}
	p.key = key
	return p
}

//Sign ..
func (p *PrivateKey) Sign(digest Sha256) Signature {
	// unsigned int buf_len = ECDSA_size(my->_key);
	// //        fprintf( stderr, "%d  %d\n", buf_len, sizeof(sha256) );
	// 		signature sig;
	// 		FC_ASSERT( buf_len == sizeof(sig) );

	// 		if( !ECDSA_sign( 0,
	// 					(const unsigned char*)&digest, sizeof(digest),
	// 					(unsigned char*)&sig, &buf_len, my->_key ) )
	// 		{
	// 			FC_THROW_EXCEPTION( exception, "signing error" );
	// 		}
	sig := new(Signature)
	return *sig
}
