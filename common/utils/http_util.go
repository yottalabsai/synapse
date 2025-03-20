package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func ProxiedClientFromEnv() *http.Client {
	transport := &http.Transport{}
	transport.Proxy = http.ProxyFromEnvironment
	httpCli := &http.Client{
		Transport: transport,
	}
	return httpCli
}

// Request sends an HTTP request. SetResult will only serialize if the response code is 200 - 299.
// errResp will serialize if the response code is not 200 - 299.
func Request[R any](
	ctx context.Context,
	client *resty.Client,
	method, url string,
	headers map[string]string,
	queryParams url.Values,
	body any,
) (data R, errResp *Response[any], err error) {
	client.JSONUnmarshal = jsoniter.Unmarshal
	result := new(Response[R])
	errResult := new(Response[any])
	resp, err := client.R().
		SetContext(ctx).
		SetQueryParamsFromValues(queryParams).
		SetHeaders(headers).
		SetBody(body).
		SetResult(result).
		SetError(errResult).
		Execute(method, url)
	// log.Log.Debugf("Request: method: %v, url: %v, param: %v, body: %+v, err: %v, httpCode: %v, response: %+v, errResp: %+v", method, url, queryParams, body, err, resp.StatusCode(), result.Data, errResult)
	if err != nil {
		return data, errResult, errors.WithMessagef(err, "resty execute request failed, method: %v, url: %v", method, url)
	}
	defer resp.RawBody().Close()

	// The service directly returns a 404 or 500 HTTP code, and the HTTP body is empty
	if resp.StatusCode() != http.StatusOK && errResult.Code == 0 {
		return data, nil, fmt.Errorf("request failed with status code: %v", resp.StatusCode())
	}

	return result.Data, errResult, nil
}
