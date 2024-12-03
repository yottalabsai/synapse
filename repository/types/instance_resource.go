package types

import (
	"time"
)

type InstanceResource struct {
	ID              int64     `gorm:"primarykey" json:"id"`
	ExternalId      string    `json:"externalId" gorm:"unique;not null"`
	InstanceType    string    `json:"instanceType" gorm:"not null"`
	ImageID         string    `json:"imageId" gorm:"not null"`
	DeviceName      string    `json:"deviceName" gorm:"not null"`
	VolumeSize      int32     `json:"volumeSize" gorm:"not null"`
	KeyPairId       string    `json:"keyPairId" gorm:"not null"`
	SecurityGroupId string    `json:"securityGroupId" gorm:"not null"`
	PrivateKey      string    `json:"privateKey"`
	PublicKey       string    `json:"publicKey"`
	InstanceId      string    `json:"instanceId" gorm:"unique;not null"`
	AssociationId   string    `json:"associationId"`
	AllocationId    string    `json:"allocationId"`
	IPAddress       string    `json:"ipAddress"`
	DomainName      string    `json:"domainName"`
	Status          int       `json:"status"`
	CreatedAt       time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (*InstanceResource) TableName() string {
	return "instance_resources"
}
