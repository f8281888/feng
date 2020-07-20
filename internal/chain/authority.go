package chain

//KeyWeight ..
type KeyWeight struct {
	permission PermissionLevel
	weight     WeightType
}

// friend bool operator == ( const key_weight& lhs, const key_weight& rhs ) {
// 	return tie( lhs.key, lhs.weight ) == tie( rhs.key, rhs.weight );
//  }

//IsEqual ..
func (k KeyWeight) IsEqual(rhs KeyWeight) bool {
	return k.permission == rhs.permission && k.weight == rhs.weight
}
