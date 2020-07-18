package chain

//ForkDatabase ForkDatabase
type ForkDatabase struct {
	ForkDatabaseImpl
}

//ForkIndex ..
type ForkIndex struct {
	ID   BlockIDType
	Prev BlockIDType
}

//如何模拟多容器索引？看来是没办法实现了，后面数据库用kv数据库来模拟吧

//ForkMultiIndexType 如果模拟不了，只能用传统的查询遍历来实现，估计效率会比较低
type ForkMultiIndexType map[*BlockIDType]*BlockState

//ForkDatabaseImpl ..
type ForkDatabaseImpl struct {
	self    *ForkDatabaseImpl
	index   *ForkMultiIndexType
	root    *BlockState
	head    *BlockState
	DataDir string
	Index   *ForkMultiIndexType
}

//Head ..
func (f ForkDatabase) Head() *BlockState {
	return f.head
}

//Root ..
func (f ForkDatabase) Root() *BlockState {
	return f.Root()
}

//RollbackHeadToRoot ..
func (f *ForkDatabase) RollbackHeadToRoot() {
	f.head = f.root
}

//Reset ..
func (f *ForkDatabase) Reset(rootBhs *BlockState) {
	f.root = &BlockState{}
	f.root.validated = true
	f.root = rootBhs
	f.head = f.root
}

//GetBlock ..
func (f *ForkDatabase) GetBlock(id BlockIDType) *BlockState {
	// err, itr := f.Index[id]
	// if err == nil {
	// 	return itr
	// }

	return nil
}
