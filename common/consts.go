package common

import (
	"fmt"
)

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
