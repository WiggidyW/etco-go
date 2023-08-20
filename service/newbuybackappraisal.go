package service

import (
	"context"

	"github.com/WiggidyW/weve-esi/proto"
)

func (s *Service) NewBuybackAppraisal(
	ctx context.Context,
	req *proto.NewBuybackAppraisalRequest,
)
