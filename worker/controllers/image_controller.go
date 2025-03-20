package controllers

import (
	"github.com/gin-gonic/gin"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"synapse/common"
	"synapse/connector/rpc"
	service2 "synapse/connector/service"
	"synapse/log"
	"synapse/utils"
	"synapse/worker/types"
	"time"
)

type ImageController struct {
	server *service2.SynapseServer
}

func NewImageController(server *service2.SynapseServer) *ImageController {
	return &ImageController{server: server}
}

func (ctl *ImageController) Render(ctx *gin.Context) {
	endpointId := ctx.Param("endpointId")

	if endpointId == "" {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}

	req := types.TextToImageRequest{}
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

	// 执行inference
	ctl.DoRender(ctx, &req)
}

func (ctl *ImageController) DoRender(ctx *gin.Context, req *types.TextToImageRequest) {
	// filter ready client
	requestID := utils.GenerateRequestId()
	flag := false
	for clientID := range service2.GlobalStreamManager.GetStreams() {
		streamDetail := service2.GlobalStreamManager.GetStreams()[clientID]
		log.Log.Infow("[search] clients", zap.Any("clientInfo", streamDetail))
		if streamDetail.Ready && streamDetail.Model == req.Model {
			// create inference request message
			msg := &synapseGrpc.Message{
				// todo: leo
			}
			if err := service2.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
				log.Log.Errorw("send message to client failed", zap.Error(err))
			} else {
				flag = true
				break
			}
		}
	}

	if !flag {
		ctx.JSON(common.HttpOk, common.ErrNoReadyClient)
		return
	}

	respChannel := rpc.GlobalChannelManager.CreateChannel(requestID)
	defer rpc.GlobalChannelManager.RemoveChannel(requestID)

	select {
	case msg := <-respChannel.InferenceResultChan:
		log.Log.Infow("[search] client response", zap.Any("msg", msg))
		response := &types.TextToImageResponse{
			Created: time.Now().Unix(),
			// todo: leo
			Data: nil,
		}

		ctx.JSON(common.HttpOk, response)
	case <-time.After(30 * time.Second):
		ctx.JSON(common.HttpOk, common.ErrTimeout)
	}

}
