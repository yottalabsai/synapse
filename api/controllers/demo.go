package controllers

import (
	"github.com/gin-gonic/gin"
	"synapse/api/types"
	"synapse/common"
	"synapse/log"
	"synapse/service"
)

type DemoController struct {
	svc *service.SwapService
}

func NewDemo(svc *service.SwapService) *DemoController {
	return &DemoController{svc: svc}
}

func (ctl *DemoController) Demo(ctx *gin.Context) {
	req := types.FindRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}
	if err := common.Validate(&req); err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("validate failed: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	res, err := ctl.svc.FindById(ctx, &req)

	if err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("save token  failed: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	common.JSON(ctx, common.HttpOk, common.Ok(res))
}
