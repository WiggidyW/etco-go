package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBContractItemsParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	ContractId        int32
}
