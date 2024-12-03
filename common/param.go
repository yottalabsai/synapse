package common

import (
	"fmt"
)

const (
	Uni4goHttpOK      int = 0
	HttpOk            int = 200
	HttpInternalError int = 500
)

var (
	ErrBadArgument          = NewApiError(107000, "Bad argument")
	ErrUnauthorized         = NewApiError(107001, "Unauthorized")
	ErrSymbolConfigNotFound = NewApiError(107002, "Symbol config not found")
	ErrAssetDelisting       = NewApiError(107028, "Trading unavailable due to the delisting of {{.Symbol}}")

	//  系统暂停的错误代码
	ErrSystemPaused = &ApiResponse{Code: -100000, Msg: "This service is temporarily unavailable. Please contact customer support."} // 系统暂停, data需要有值

	//  后管
	ErrSystemBusy = NewApiError(107020, "System busy")
)

// 钱包服务错误码
var (
	ServiceCodeTokenTrades    = 107029 // 币币交易服务
	ServiceCodeTokenTransfers = 107030 // 划转服务
)

func ConvertWalletErrCode(code int, msg string) *ApiError {
	return &ApiError{
		Code: ApiErrorCode(code),
		Msg:  msg,
	}
}

type ApiErrorCode int

type ApiError struct {
	Code         ApiErrorCode   `json:"code"`
	Msg          string         `json:"msg"`
	TemplateData map[string]any `json:"-"`
}

func NewApiError(code ApiErrorCode, msg string) *ApiError {
	return &ApiError{
		Code: code,
		Msg:  msg,
	}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("code: %v, msg: %v", e.Code, e.Msg)
}

type Validatable interface {
	// Validate 验证请求参数
	Validate() error
}

func Validate(r any) error {
	if validatable, ok := r.(Validatable); ok {
		return validatable.Validate()
	}
	return nil
}

type ApiResponse struct {
	Code         int            `json:"code"`
	Msg          string         `json:"msg"`
	Data         any            `json:"data"`
	TemplateData map[string]any `json:"-"`
}

func Ok(data any) *ApiResponse {
	return &ApiResponse{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func GuessError(err error, unknownErrHandler func(error)) (httpCode int, apiResp *ApiResponse) {
	if apiErr, ok := err.(*ApiError); ok {
		return HttpOk, &ApiResponse{
			Code:         int(apiErr.Code),
			Msg:          apiErr.Error(),
			TemplateData: apiErr.TemplateData,
		}
	}
	unknownErrHandler(err)
	return HttpInternalError, internalError()
}

func internalError() *ApiResponse {
	return &ApiResponse{
		Code: 500,
		Msg:  "Internal error",
	}
}
