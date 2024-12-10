package types

import (
	"time"
)

type ServerlessResource struct {
	ID         int64     `gorm:"primarykey" json:"id"`
	EndpointId string    `json:"endpointId" gorm:"unique;not null"`
	Model      string    `json:"model" gorm:"not null"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (*ServerlessResource) TableName() string {
	return "serverless_resources"
}
