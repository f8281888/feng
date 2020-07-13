package chain

//PermissionLevel ..
type PermissionLevel struct {
	Actor      AccountName
	Permission PermissionName
}

//Action ..
type Action struct {
	Account       AccountName
	Name          ActionName
	Authorization []PermissionLevel
	data          []byte
}
