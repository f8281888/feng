package chain

//Config ..
type Config struct {
	maxBlockNetUsage               uint64
	targetBlockNetUsagePct         uint32
	maxTransactionNetUsage         uint32
	basePerTransactionNetUsage     uint32
	netUsageLeeway                 uint32
	contextFreeDiscountNetUsageNum uint32
	contextFreeDiscountNetUsageDen uint32
	maxBlockCPUUsage               uint32
	targetBlockCPUUsagePct         uint32
	maxTransactionCPUUsage         uint32
	minTransactionCPUUsage         uint32
	maxTransactionLifetime         uint32
	deferredTrxExpirationWindow    uint32
	maxTransactionDelay            uint32
	maxInlineActionSize            uint32
	maxInlineActionDepth           uint16
	maxAuthorityDepth              uint16
}
