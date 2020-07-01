package chaininterface

import (
	"feng/internal/chain"
)

//PreAcceptedBlock ..
type PreAcceptedBlock struct {
	signedBlock *chain.SignedBlock
}

//RejectedBlock ..
type RejectedBlock struct {
	signedBlock *chain.SignedBlock
}

//AcceptedBlock ..
type AcceptedBlock struct {
	signedBlock *chain.BlockState
}

//IrreversibleBlock ..
type IrreversibleBlock struct {
	signedBlock *chain.BlockState
}

//AcceptedTransaction ..
type AcceptedTransaction struct {
}

//AppliedTransaction ..
type AppliedTransaction struct {
}
