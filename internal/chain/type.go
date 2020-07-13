package chain

import (
	"feng/internal/fc/crypto"
)

//ActionName ..
type ActionName = Name

//ScopeName ..
type ScopeName = Name

//AccountName ..
type AccountName = Name

//PermissionName ..
type PermissionName = Name

//TableName ..
type TableName = Name

//BlockIDType ..
type BlockIDType = crypto.Sha256

//ChecksumType ..
type ChecksumType = crypto.Sha256

//Checksum256Type ..
type Checksum256Type = crypto.Sha256

//TransactionIDType ..
type TransactionIDType = ChecksumType

//DigestType ..
type DigestType = ChecksumType

//WeightType ..
type WeightType = uint16

//BlockNumType ..
type BlockNumType = uint32

//ShareType ..
type ShareType = int64

//PublicKeyType ..
type PublicKeyType = crypto.PublicKey

//PrivateKeType ..
type PrivateKeType = crypto.PrivateKey

//SignatureType ..
type SignatureType = crypto.Signature
