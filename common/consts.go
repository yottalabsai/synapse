package common

import (
	"fmt"
)

const StatusStopped = -1
const StatusInit = 0
const StatusSuccess = 1

const ResourceStatusCheckLockPrefix = "lock:resource:check"
const JobResourceStatusCheckSecond = 10

var JobResourceStatusCheckSpec = fmt.Sprintf("@every %ds", JobResourceStatusCheckSecond)

const InstanceRunningCheckLokPrefix = "lock:instance:running:check"
const InstanceStoppingCheckLockPrefix = "lock:instance:stopping:check"
const InstanceTerminatingCheckLockPrefix = "lock:instance:terminating:check"
const JobInstanceStatusCheckSecond = 10

var JobInstanceStatusCheckSpec = fmt.Sprintf("@every %ds", JobInstanceStatusCheckSecond)
