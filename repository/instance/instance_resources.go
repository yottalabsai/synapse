package instance

import (
	"gorm.io/gorm"
	repo "synapse/repository"
	"synapse/repository/types"
)

type ManageInstanceRepo struct {
	repo.BaseRepo[*types.InstanceResource]
}

func NewInstanceRepository(db *gorm.DB) *ManageInstanceRepo {
	return &ManageInstanceRepo{repo.BaseRepo[*types.InstanceResource]{DB: db}}
}

func (b *ManageInstanceRepo) FindById(id uint64) (*types.InstanceResource, error) {
	record := new(types.InstanceResource)
	if err := b.DB.Where("id = ?", id).First(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
