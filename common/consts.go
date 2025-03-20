package common

import (
	"fmt"
)

const ServiceConnector = "connector"
const ServiceWorker = "worker"

const StatusStopped = -1
const StatusInit = 0
const StatusSuccess = 1

const InferencePublicListCheckLockPrefix = "lock:inference:public:list:check"
const JobResourceStatusCheckSecond = 10

var JobInferencePublicListCheckSpec = fmt.Sprintf("@every %ds", JobResourceStatusCheckSecond)

const (
	ServiceYottaSaaS = "yotta-saas"

	UrlPathInferencePublicList = "/api/inference/public/list"
)

type ModelType string

const (
	Inference   ModelType = "1"
	TextToImage ModelType = "2"
)

func (m ModelType) ToString() string {
	if m == Inference {
		return "1:Inference"
	} else {
		return "2:TextToImage"
	}
}
