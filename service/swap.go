package service

import (
	"context"
	"gorm.io/gorm"
	"synapse/api/types"
	"synapse/repository/instance"
)

type SwapService struct {
	db *gorm.DB
}

func NewSwapService(db *gorm.DB) *SwapService {
	return &SwapService{
		db: db,
	}
}

func (svc *SwapService) FindById(ctx context.Context, req *types.FindRequest) (*types.FindResponse, error) {

	repo := instance.NewInstanceRepository(svc.db.WithContext(ctx))

	result, err := repo.FindById(req.InstanceId)
	if err != nil {
		return nil, err
	}

	return &types.FindResponse{
		ID:           result.ID,
		ExternalId:   result.ExternalId,
		InstanceType: result.InstanceType,
	}, nil
}
