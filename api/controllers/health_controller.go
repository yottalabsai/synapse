package controllers

import (
	"github.com/gin-gonic/gin"
	"synapse/common"
)

func Health(c *gin.Context) {
	common.JSON(c, common.HttpOk, common.Ok("success"))
}
