package common

//ScopeExit ..
type ScopeExit struct {
	callback func()
	caceled  bool
}

//MakeScopeExit ..
func MakeScopeExit(c func()) ScopeExit {
	a := ScopeExit{callback: c}
	return a
}
