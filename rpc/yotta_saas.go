package rpc

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/url"
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
	ModelID          string `json:"modelId"`
	ModelDisplayName string `json:"modelDisplayName"`
	ModelName        string `json:"modelName"`
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

	return res, nil
}
