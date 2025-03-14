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

type InferenceController struct {
	server *service.SynapseServer
}

func NewInferenceController(server *service.SynapseServer) *InferenceController {
	return &InferenceController{server: server}
}

func (ctl *InferenceController) Inference(ctx *gin.Context) {
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

	// run inference
	ctl.DoInference(ctx, &req)
}

func (ctl *InferenceController) DoInference(ctx *gin.Context, req *types.InferenceMessageRequest) {

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
		log.Log.Infow("[search] clients", zap.Any("clientInfo", streamDetail))

		inferenceMessage := &synapseGrpc.InferenceMessage{
			Temperature:       req.Temperature,
			TopP:              req.TopP,
			MaxTokens:         req.MaxTokens,
			FrequencyPenalty:  req.FrequencyPenalty,
			PresencePenalty:   req.PresencePenalty,
			RepetitionPenalty: req.RepetitionPenalty,
			Model:             req.Model,
			Stream:            req.Stream,
			Messages:          messages,
		}

		if req.Stream {
			inferenceMessage.StreamOptions = &synapseGrpc.StreamOptions{
				IncludeUsage: req.StreamOptions.IncludeUsage,
			}
		}

		if streamDetail.Ready && streamDetail.Model == req.Model {
			// create inference request message
			msg := &synapseGrpc.YottaLabsStream{
				MessageId: requestID,
				Timestamp: time.Now().Unix(),
				ClientId:  clientID,
				Payload: &synapseGrpc.YottaLabsStream_InferenceMessage{
					InferenceMessage: inferenceMessage,
				},
			}
			if err := service.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
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

	respChannel := service.GlobalChannelManager.CreateChannel(requestID)
	defer service.GlobalChannelManager.RemoveChannel(requestID)

	// config SSE header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	ctx.Stream(func(w io.Writer) bool {
		// return true: continue streaming
		// return false: end streaming
		select {
		case result := <-respChannel.InferenceResultChan:
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
				return false
			}
			return true
		case <-ctx.Request.Context().Done():
			ctx.SSEvent("error", "client disconnected")
			return false // stop streaming
		case <-time.After(30 * time.Second):
			// stop streaming
			ctx.SSEvent("error", "timeout")
			return false // timeout
		}
	})
}
