package types

type FindRequest struct {
	InstanceId uint64 `json:"instanceId,string"  form:"instanceId"  binding:"required"`
}

type FindResponse struct {
	ID           int64  `json:"id"`
	ExternalId   string `json:"externalId"`
	InstanceType string `json:"instanceType"`
}
