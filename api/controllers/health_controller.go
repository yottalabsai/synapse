package controllers

import (
	"github.com/gin-gonic/gin"
	"synapse/common"
	"synapse/connector/service"
)

type HealthController struct {
	statusService *service.StatusService
}

func NewHealthController(server *service.StatusService) *HealthController {
	return &HealthController{statusService: server}
}

func (h *HealthController) Health(c *gin.Context) {
	common.JSON(c, common.HttpOk, common.Ok("success"))
}

func (h *HealthController) Status(c *gin.Context) {
	common.JSON(c, common.HttpOk, common.Ok(h.statusService.GetStatus()))
}
