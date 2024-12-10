package controllers

import (
	"github.com/gin-gonic/gin"
	"synapse/api/types"
	"synapse/common"
	"synapse/log"
	"synapse/service"
)

type ServerlessController struct {
	svc *service.ServerlessService
}

func NewServerlessController(svc *service.ServerlessService) *ServerlessController {
	return &ServerlessController{svc: svc}
}

func (ctl *ServerlessController) FindByEndpointId(ctx *gin.Context) {
	endpointId := ctx.Param("endpointId")

	if endpointId == "" {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}

	res, err := ctl.svc.FindByEndpointId(ctx, endpointId)

	if err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("find serverless error: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	common.JSON(ctx, common.HttpOk, common.Ok(res))
}

func (ctl *ServerlessController) CreateEndpoint(ctx *gin.Context) {
	req := types.CreateServerlessResourceRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	res, err := ctl.svc.Create(ctx, &req)

	if err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("create serverless error: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	common.JSON(ctx, common.HttpOk, common.Ok(res))
}

func (ctl *ServerlessController) Inference(ctx *gin.Context) {
	endpointId := ctx.Param("endpointId")

	if endpointId == "" {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}

	req := types.InferenceMessageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	res, err := ctl.svc.Inference(ctx, endpointId, &req)

	if err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("inference error: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	common.JSON(ctx, common.HttpOk, common.Ok(res))
}
