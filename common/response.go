package common

import (
	"github.com/gin-gonic/gin"
)

// JSON 封装了ctx.JSON，对msg进行国际化处理
// 任何响应都应该通过调用此方法进行
// Body必须是以下类型:
// - *ApiError
// - *ApiResponse
// 其他类型将导致panic
func JSON(ctx *gin.Context, httpCode int, body any) {
	ctx.JSON(httpCode, body)
}
