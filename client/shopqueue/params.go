package shopqueue

import (
	"github.com/WiggidyW/chanresult"
)

type ShopQueueParams struct {
	// will send nothing if not modified
	// will send struct{}{} if modified successful
	// will send error if not
	ChnSendModifyDone *chanresult.ChanSendResult[struct{}]
}
