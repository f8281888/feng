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

// namespace incoming {
// 	namespace channels {
// 	   using block                 = channel_decl<struct block_tag, signed_block_ptr>;
// 	   using transaction           = channel_decl<struct transaction_tag, packed_transaction_ptr>;
// 	}

// 	namespace methods {
// 	   // synchronously push a block/trx to a single provider
// 	   using block_sync            = method_decl<chain_plugin_interface, bool(const signed_block_ptr&, const std::optional<block_id_type>&), first_provider_policy>;
// 	   using transaction_async     = method_decl<chain_plugin_interface, void(const packed_transaction_ptr&, bool, next_function<transaction_trace_ptr>), first_provider_policy>;
// 	}
//  }
