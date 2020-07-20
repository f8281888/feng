package chain

//BlockIntervalMs ..
var BlockIntervalMs int = 500

//BlockIntervalUs ..
var BlockIntervalUs int = BlockIntervalMs * 1000

//Percent100 ..
var Percent100 int = 10000

//Percent1 ..
var Percent1 int = 100

//DefaultSigCPUBillPct ..
var DefaultSigCPUBillPct uint32 = uint32(50 * Percent1)

//DefaultBlockCPUEffortPct ..
var DefaultBlockCPUEffortPct uint32 = uint32(80 * Percent1)

//FengPrercent ..
func FengPrercent(value uint64, percentage uint32) uint64 {
	return value * uint64(percentage) / uint64(Percent100)
}

//DefaultSubjectiveCPULeewayUs ..
var DefaultSubjectiveCPULeewayUs uint32 = 31000

//DefaultControllerThreadPoolSize ..
var DefaultControllerThreadPoolSize int16 = 2

//ProducerRepetitions ..
var ProducerRepetitions uint32 = 12
