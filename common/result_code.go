package common

import "fmt"

type ResultCode struct {
	Code    int
	Message string
}

var (
	SUCCESS                 = ResultCode{200, "success"}
	PARAMETER_ERROR         = ResultCode{400, "parameter error"}
	SYSTEM_ERROR            = ResultCode{500, "system error"}
	LIMIT                   = ResultCode{503, "Too Frequent Request"}
	RESOURCE_NOT_FOUND      = ResultCode{100000, "Resource not found"}
	START_SERVICE_FAILED    = ResultCode{100001, "Cannot start service"}
	STOP_SERVICE_FAILED     = ResultCode{100002, "Cannot stop service"}
	STOPPED_SERVICE_ALREADY = ResultCode{100003, "Stop service already"}
)

// Error implements the error interface for CustomError
func (e *ResultCode) Error() string {
	return fmt.Sprintf("Message: %s", e.Message)
}

// NewBusinessError creates a new CustomError with the given ResultCode
func NewBusinessError(e ResultCode) error {
	return &e
}
