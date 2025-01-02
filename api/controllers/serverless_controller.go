package controllers

import (
	"github.com/gin-gonic/gin"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"io"
	"strings"
	"synapse/api/types"
	"synapse/common"
	"synapse/log"
	"synapse/service"
	"synapse/utils"
	"time"
)

type ServerlessController struct {
	svc    *service.ServerlessService
	server *service.SynapseServer
}

func NewServerlessController(svc *service.ServerlessService, server *service.SynapseServer) *ServerlessController {
	return &ServerlessController{svc: svc, server: server}
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

	// 执行inference
	ctl.DoInference(ctx, &req)
}

func (ctl *ServerlessController) DoInference(ctx *gin.Context, req *types.InferenceMessageRequest) {

	messages := make([]*synapseGrpc.Message, len(req.Messages))
	index := 0
	for _, message := range req.Messages {
		messages[index] = &synapseGrpc.Message{
			Content: message.Content,
			Role:    message.Role,
		}
	}

	// filter ready client
	requestID := utils.GenerateRequestId()
	flag := false
	for clientID := range service.GlobalStreamManager.GetStreams() {
		streamDetail := service.GlobalStreamManager.GetStreams()[clientID]
		if streamDetail.Ready && streamDetail.Model == req.Model {
			// create inference request message
			msg := &synapseGrpc.YottaLabsStream{
				MessageId: requestID,
				Timestamp: time.Now().Unix(),
				ClientId:  clientID,
				Payload: &synapseGrpc.YottaLabsStream_InferenceMessage{
					InferenceMessage: &synapseGrpc.InferenceMessage{
						Temperature:       req.Temperature,
						TopP:              req.TopP,
						MaxTokens:         req.MaxTokens,
						FrequencyPenalty:  req.FrequencyPenalty,
						PresencePenalty:   req.PresencePenalty,
						RepetitionPenalty: req.RepetitionPenalty,
						Model:             req.Model,
						Stream:            req.Stream,
						StreamOptions: &synapseGrpc.StreamOptions{
							IncludeUsage: req.StreamOptions.IncludeUsage,
						},
						Messages: messages,
					},
				},
			}
			if err := service.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
				log.Log.Error("send message to client failed", zap.Error(err))
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

	respChannel := service.GlobalChannelManager.CreateChannel(requestID)
	defer service.GlobalChannelManager.RemoveChannel(requestID)

	// 设置 SSE 相关的 header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	ctx.Stream(func(w io.Writer) bool {
		// return true: continue streaming
		// return false: end streaming
		select {
		case result := <-respChannel.ResultChan:
			// remove data: prefix
			content := result.InferenceResult.Content
			if len(content) > 6 && content[:6] == "data: " {
				content = content[6:]
			}
			// send data to client, no need event
			ctx.SSEvent("", " "+content)
			// check if inference is done
			if strings.Index(content, "[DONE]") == 0 {
				// stop streaming
				service.GlobalChannelManager.RemoveChannel(requestID)
				return false
			}
			return true
		case <-ctx.Request.Context().Done():
			ctx.SSEvent("error", "client disconnected")
			return false // stop streaming
		case <-time.After(30 * time.Second):
			// stop streaming
			service.GlobalChannelManager.RemoveChannel(requestID)
			ctx.SSEvent("error", "timeout")
			return false // timeout
		}
	})
}
