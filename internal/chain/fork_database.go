package chain

//ForkDatabase ForkDatabase
type ForkDatabase struct {
	myForkDatabase *ForkDatabaseImpl
}

//ForkDatabaseImpl ..
type ForkDatabaseImpl struct {
	self *ForkDatabaseImpl
	//fork_multi_index_type
	root    *BlockState
	head    *BlockState
	DataDir string
}

//Head ..
func (f ForkDatabase) Head() *BlockState {
	return f.myForkDatabase.head
}

//Root ..
func (f ForkDatabase) Root() *BlockState {
	return f.myForkDatabase.root
}

//RollbackHeadToRoot ..
func (f *ForkDatabase) RollbackHeadToRoot() {
	f.myForkDatabase.head = f.myForkDatabase.root
}

//Reset ..
func (f *ForkDatabase) Reset(rootBhs *BlockState) {
	f.myForkDatabase.root = &BlockState{}
	f.myForkDatabase.root.validated = true
	f.myForkDatabase.root = rootBhs
	f.myForkDatabase.head = f.myForkDatabase.root
}

//Set ..
func (f *ForkDatabase) Set(fork *ForkDatabaseImpl) {
	f.myForkDatabase = fork
}
