package chain

import (
	"feng/internal/log"
)

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
type ForkMultiIndexType map[*BlockIDType]BlockState

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

//Add ..
func (f *ForkDatabase) Add(n *BlockState, ignoreDuplicate bool) {
	f.add(n, ignoreDuplicate, false, func(BlockTimestamp, []DigestType, []DigestType) {})
}

func (f *ForkDatabase) add(n *BlockState, ignoreDuplicate bool, validate bool, validator func(BlockTimestamp, []DigestType, []DigestType)) {
	if f.root == nil {
		log.Assert("root not yet set")
	}

	if n == nil {
		log.Assert("attempt to add null block state")
	}

	prevBh := f.getBlockHeader(&n.header.previous)
	if prevBh == nil {
		log.Assert("unlinkable block id :%s, previous:%s", n.id.String(), n.header.previous.String())
	}

	if validate {
		exts := n.headerExts
		//exts.count(protocol_feature_activation::extension_id()) > 0 第一个元素出现的次数
		if exts[0] != (BlockHeaderExtensionTypes{}) {
			newProtocolFeatures := exts[0].protocolFeatureActivation.ProtocolFeatures
			validator(n.header.Timestamp, prevBh.activatedProtocolFeatures.ProtocolFeatures, newProtocolFeatures)
		}
	}

	_, ok := (*f.index)[n.id]
	if !ok {
		(*f.index)[n.id] = *n
	}

	for _, i := range *f.index {
		if i.IsValid() {
			f.head = &i
		}
	}
}

func (f *ForkDatabase) getBlockHeader(id *BlockIDType) *BlockState {
	byIDIdx, ok := (*f.index)[id]
	if f.root.id == id {
		return f.root
	}

	if ok {
		return &byIDIdx
	}

	return &BlockState{}
}

func (f ForkDatabase) pendingHead() *BlockState {
	for _, i := range *f.index {
		if !i.IsValid() {
			if f.firstPreferred(i, *f.head) {
				return &i
			}
		}
	}

	return f.head
}

func (f ForkDatabase) firstPreferred(lhs BlockState, rhs BlockState) bool {
	return lhs.dposIrreversibleBlocknum > rhs.dposIrreversibleBlocknum && lhs.BlockNum > rhs.BlockNum
}
