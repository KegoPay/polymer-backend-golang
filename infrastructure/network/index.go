package network

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/logger/metrics"
)

type NetworkController struct {
	BaseUrl    string
	HttpClient *http.Client
}

func (network *NetworkController) InitialiseNetworkClient() {
	network.HttpClient = &http.Client{}
	network.HttpClient.Transport = metrics.MetricMonitor.GetRoundTripper(context.Background())
}

func (network *NetworkController) Get(path string, headers *map[string]string, params *map[string]string) (response *[]byte, statusCode *int, err error) {
	if network.HttpClient == nil {
		network.InitialiseNetworkClient()
	}

	req, err := http.NewRequest("GET", network.BaseUrl+path, nil)
	if err != nil {
		logger.Error(errors.New("could not initiate GET request"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	network.setHeaders(headers, req)
	network.setParams(params, req)

	res, err := network.HttpClient.Do(req)
	if err != nil {
		logger.Error(errors.New("could not complete GET request"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(errors.New("could not read GET request body"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	return &resBody, &res.StatusCode, nil
}

func (network *NetworkController) Post(path string, headers *map[string]string, body any, params *map[string]string) (response *[]byte, statusCode *int, err error) {
	if network.HttpClient == nil {
		network.InitialiseNetworkClient()
	}

	parsedBody, err := json.Marshal(body)
	if err != nil {
		logger.Error(errors.New("error converting body to JSON"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", network.BaseUrl+path, bytes.NewBuffer(parsedBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logger.Error(errors.New("could not initiate POST request"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	network.setHeaders(headers, req)
	network.setParams(params, req)

	defer req.Body.Close()

	res, err := network.HttpClient.Do(req)
	if err != nil {
		logger.Error(errors.New("could not complete POST request"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return &resBody, &res.StatusCode, nil
}

func (network *NetworkController) setHeaders(headers *map[string]string, req *http.Request) {
	if headers == nil {
		return
	}

	for k := range *headers {
		req.Header.Add(k, (*headers)[k])
	}
}

func (network *NetworkController) setParams(params *map[string]string, req *http.Request) {
	if params == nil {
		return
	}

	q := req.URL.Query()
	for k := range *params {
		q.Add(k, (*params)[k])
	}
	req.URL.RawQuery = q.Encode()
}
