package common

import (
	"fmt"
)

const KeyPairInactive = 0
const KeyPairActive = 1

const S3MountPath = "/data"

const ServiceStatusStopped = -1
const ServiceStatusInit = 0
const ServiceStatusSuccess = 1

const InstanceStatusFailed = -1
const InstanceStatusInit = 0
const InstanceStatusRunning = 1
const InstanceStatusStopping = 2
const InstanceStatusStopped = 3
const InstanceStatusShuttingDown = 4
const InstanceStatusTerminated = 5

const ResourceStatusCheckLockPrefix = "lock:resource:check"
const JobResourceStatusCheckSecond = 10

var JobResourceStatusCheckSpec = fmt.Sprintf("@every %ds", JobResourceStatusCheckSecond)

const InstanceRunningCheckLokPrefix = "lock:instance:running:check"
const InstanceStoppingCheckLockPrefix = "lock:instance:stopping:check"
const InstanceTerminatingCheckLockPrefix = "lock:instance:terminating:check"
const JobInstanceStatusCheckSecond = 10

var JobInstanceStatusCheckSpec = fmt.Sprintf("@every %ds", JobInstanceStatusCheckSecond)
