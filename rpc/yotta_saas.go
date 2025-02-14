package rpc

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"synapse/common"
	"synapse/config"
	"synapse/utils"
)

type YottaSaaSClient struct {
	cfg   *config.ServiceConfig
	resty *resty.Client
}

func NewYottaSaaSClient(cfg *config.ServiceConfig, client *resty.Client) *YottaSaaSClient {
	serviceConfig := config.MustGetServiceConfig(common.ServiceYottaSaaS)[0]
	return &YottaSaaSClient{
		cfg:   &serviceConfig,
		resty: client,
	}
}

type ModelInfo struct {
	ModelID          string           `json:"modelId"`
	ModelDisplayName string           `json:"modelDisplayName"`
	ModelName        string           `json:"modelName"`
	ModelType        common.ModelType `json:"modelType"`
	Ready            bool
}

func (c *YottaSaaSClient) FindInferencePublicList(ctx context.Context) (*[]ModelInfo, error) {
	params := url.Values{}
	res, _, err := utils.Request[*[]ModelInfo](
		ctx,
		c.resty,
		resty.MethodGet,
		lo.Must1(url.JoinPath(c.cfg.Endpoint, common.UrlPathInferencePublicList)),
		c.cfg.Headers,
		params,
		nil,
	)

	if err != nil {
		return nil, errors.WithMessagef(err, "get public model list failed")
	}

	for i := range *res {
		// 判断是否包含 mit-han-lab/svdq-int4-flux.1-schnell 和 black-forest-labs/FLUX.1-schnell
		// Efficient-Large-Model/Sana_1600M_1024px_BF16_diffusers 和 mit-han-lab/svdq-int4-sana-1600m
		// black-forest-labs/FLUX.1-dev 和 mit-han-lab/svdq-int4-flux.1-dev
		if strings.Contains((*res)[i].ModelName, "mit-han-lab/svdq-int4-flux.1-schnell") ||
			strings.Contains((*res)[i].ModelName, "black-forest-labs/FLUX.1-schnell") ||
			strings.Contains((*res)[i].ModelName, "mit-han-lab/svdq-int4-sana-1600m") ||
			strings.Contains((*res)[i].ModelName, "Efficient-Large-Model/Sana_1600M_1024px_BF16_diffusers") ||
			strings.Contains((*res)[i].ModelName, "black-forest-labs/FLUX.1-dev") ||
			strings.Contains((*res)[i].ModelName, "mit-han-lab/svdq-int4-flux.1-dev") {
			(*res)[i].ModelType = common.TextToImage
		} else {
			(*res)[i].ModelType = common.Inference
		}
	}

	return res, nil
}
