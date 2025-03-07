package common

import (
	"github.com/gin-gonic/gin"
)

// JSON wraps ctx.JSON and processes msg for internationalization
// Any response should be made by calling this method
// Body must be one of the following types:
// - *ApiError
// - *ApiResponse
// Other types will cause a panic
func JSON(ctx *gin.Context, httpCode int, body any) {
	ctx.JSON(httpCode, body)
}
