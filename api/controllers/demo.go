package controllers

import (
	"synapse/service"

	"fmt"
	"synapse/log"
	"time"

	"go.uber.org/zap"

	synapse_grpc "github.com/yottalabsai/endorphin/pkg/services/synapse"

	"io"

	"github.com/gin-gonic/gin"
)

type DemoController struct {
	svc *service.SwapService
}

func NewDemo(svc *service.SwapService) *DemoController {
	return &DemoController{svc: svc}
}

func (ctl *DemoController) Demo(ctx *gin.Context) {
	// 设置 SSE 相关的 header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

	// create inference request message
	msg := &synapse_grpc.StreamMessage{
		Base: &synapse_grpc.BaseMessage{
			MessageId: requestID,
			Timestamp: time.Now().Unix(),
			SenderId:  "Scheduler-1",
		},
		Payload: &synapse_grpc.StreamMessage_Inference{
			Inference: &synapse_grpc.InferenceRequest{
				ModelId: "model-1",
				Contents: []*synapse_grpc.InferenceContent{
					{
						Content: "input-1",
					},
				},
			},
		},
	}

	// send message to all connected agents
	for clientID := range service.GlobalStreamManager.GetStreams() {
		if err := service.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
			log.Log.Error("send message to client failed", zap.Error(err))
			return
		}
	}

	respChannel := service.GlobalRequestManager.CreateChannel(requestID)
	defer service.GlobalRequestManager.RemoveChannel(requestID)

	// use SSE to send result
	ctx.Stream(func(w io.Writer) bool {
		// return true: continue streaming
		// return false: end streaming
		select {
		case result := <-respChannel.ResultChan:
			log.Log.Info("inference response", zap.Any("result", result))
			if result.InferenceResponse.GetResult() == "done" {
				// stop streaming
				service.GlobalRequestManager.RemoveChannel(requestID)
				return false
			}
			// send data to client
			ctx.SSEvent("message", gin.H{
				"status": "success",
				"type":   "inference_response",
				"data": gin.H{
					"result": result.InferenceResponse,
				},
			})
			return true
		case <-ctx.Request.Context().Done():
			ctx.SSEvent("error", gin.H{
				"status": "error",
				"error":  "client disconnected",
			})
			return false // stop streaming
		case <-ctx.Request.Context().Done():
			// client disconnected
			return false
		}
	})
}
