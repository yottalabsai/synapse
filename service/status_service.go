package service

import "synapse/api/types"

type StatusService struct {
	manager *StreamManager
}

var modelRpmMap = map[string]int{
	"meta-llama/Llama-3.2-3B-Instruct":                       24,
	"unsloth/Meta-Llama-3.1-8B-Instruct-bnb-4bit":            20,
	"thesven/Mistral-7B-Instruct-v0.3-GPTQ":                  21,
	"mit-han-lab/svdq-int4-flux.1-schnell":                   19,
	"mit-han-lab/svdq-int4-sana-1600m":                       15,
	"black-forest-labs/FLUX.1-schnell":                       18,
	"black-forest-labs/FLUX.1-dev":                           12,
	"Efficient-Large-Model/Sana_1600M_1024px_BF16_diffusers": 11,
}

func NewStatusService(db *StreamManager) *StatusService {
	return &StatusService{manager: db}
}

func (svc *StatusService) GetStatus() types.StatusResponse {
	return types.StatusResponse{
		Resources: types.Resources{
			TotalNodes: svc.GetTotalNode(),
		},
		Models: types.Models{
			List: svc.GetModelMap(),
		},
	}
}

func (svc *StatusService) GetTotalNode() int {
	streamMap := svc.manager.GetStreams()
	return len(streamMap)
}

func (svc *StatusService) GetModelMap() []*types.ModelInfo {

	streamMap := svc.manager.GetStreams()
	modelInfos := make([]*types.ModelInfo, len(streamMap))

	for _, streamDetail := range streamMap {
		if streamDetail.Ready == true {
			_ = append(modelInfos, &types.ModelInfo{
				Model: streamDetail.Model,
				Count: 1,
				TPM:   modelRpmMap[streamDetail.Model],
			})
		}
	}

	return modelInfos
}
