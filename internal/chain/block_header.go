package chain

import (
	"container/list"
	"feng/internal/fc/common"
	"feng/internal/fc/crypto"
)

//BlockHeader 头
type BlockHeader struct {
	//时间戳 生成区块的时间
	Timestamp BlockTimestamp
	//生产区块的节点
	Producer Name
	//区块的确认数
	Confirmed uint16
	previous  crypto.Sha256
	//交易的默克尔树  区块中全部交易的默克尔树的哈希值
	transactionMroot crypto.Sha256
	//action的默克尔树 区块中全部action 的默克尔树的哈希值
	actionMroot crypto.Sha256
	//版本号 见证人排序版本号
	scheduleVersion uint32
	//区块的下一个见证人，可以为空
	newProducers ProducerScheduleType
	//扩展类型
	extensionsType map[uint16]list.List
}

//ID ..
func (b *BlockHeader) ID() crypto.Sha256 {
	a := crypto.Sha256{}
	a.New("0x0000")
	return a
}

//numFromID ..
func (b *BlockHeader) numFromID(id crypto.Sha256) uint32 {
	return uint32(common.BytesToInt(id.Hash))
}

//BlockNum ..
func (b *BlockHeader) BlockNum() uint32 {
	return b.numFromID(b.previous) + 1
}

//SignedBlockHeader ..
type SignedBlockHeader struct {
	BlockHeader
	producerSignature string
}

//ID ..
func (s SignedBlockHeader) ID() crypto.Sha256 {
	// block_id_type result = digest(); //fc::sha256::hash(*static_cast<const block_header*>(this));
	// result._hash[0] &= 0xffffffff00000000;
	// result._hash[0] += fc::endian_reverse_u32(block_num()); // store the block num in the ID, 160 bits is plenty for the hash
	// return result;
	a := crypto.Sha256{}
	a.New("0x0000")
	return a
}
