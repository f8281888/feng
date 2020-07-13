package nodeenum

const (
	//OtherFail ..
	OtherFail = -1
	//InitializeFail ..
	InitializeFail = -2
	//Success ..
	Success = 0
	//BadAlloc ..
	BadAlloc = 1
	//DatabaseDirty ..
	DatabaseDirty = 2
	//FixedRecesrsible ..
	FixedRecesrsible = 0
	//ExtractedGenesis ..
	ExtractedGenesis = 0
	//NodeManagementSuccess ..
	NodeManagementSuccess = 5
)

const (
	//None ..
	None = byte('0')
	//Producers ..
	Producers = byte('1') << 0
	//Specified ..
	Specified = byte('1') << 1
	//Any ..
	Any = byte('1') << 2
)

const (
	//Both ..
	Both = byte('0')
	//TransactionsOnly ..
	TransactionsOnly = byte('1')
	//BlocksOnly ..
	BlocksOnly = byte('2')
)

const (
	//NoReason no reason to go away
	NoReason = iota
	//Self  the connection is to itself
	Self = iota
	//Duplicate the connection is redundant
	Duplicate = iota
	//WrongChain the peer's chain id doesn't match
	WrongChain = iota
	//WrongVersion the peer's network version doesn't match
	WrongVersion = iota
	//Forked the peer's irreversible blocks are different
	Forked = iota
	//Unlinkable the peer sent a block we couldn't use
	Unlinkable = iota
	//BadTransaction the peer sent a transaction that failed verification
	BadTransaction = iota
	//Validation the peer sent a block that failed validation
	Validation = iota
	//BenignOther reasons such as a timeout. not fatal but warrant resetting
	BenignOther = iota
	//FatalOther a catch-all for errors we don't have discriminated
	FatalOther = iota
	//Authentication peer failed authenicatio
	Authentication = iota
)

const (
	//IDListModeNone ..
	IDListModeNone = iota
	//IDListModeCatchUp ..
	IDListModeCatchUp = iota
	//IDListModeLastIrrCatchUp ..
	IDListModeLastIrrCatchUp = iota
	//IDListModeNormal ..
	IDListModeNormal = iota
)

const (
	//StagesLibCatchup ..
	StagesLibCatchup = iota
	//StagesHeadCatchup ..
	StagesHeadCatchup = iota
	//StagesInSync ..
	StagesInSync = iota
)

const (
	//Porducing ..
	Porducing = iota
	//Speculating ..
	Speculating = iota
)
