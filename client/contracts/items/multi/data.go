package multi

import (
	i "github.com/WiggidyW/weve-esi/client/contracts/items"
)

type ContractItems struct {
	ContractId    int32
	ContractItems []i.ContractItem
}
