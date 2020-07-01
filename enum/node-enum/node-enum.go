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
