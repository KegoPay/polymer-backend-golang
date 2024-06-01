package network

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/logger"
)

type NetworkController struct {
	BaseUrl    string
	HttpClient *http.Client
}

func (network *NetworkController) InitialiseNetworkClient() {
	network.HttpClient = &http.Client{}
	network.HttpClient.Transport = logger.RequestMetricMonitor.GetRoundTripper(context.Background())
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

func (network *NetworkController) Post(path string, headers *map[string]string, body any, params *map[string]string, formdata bool, files *map[string][]byte) (response *[]byte, statusCode *int, err error) {
	if network.HttpClient == nil {
		network.InitialiseNetworkClient()
	}
	var buf bytes.Buffer
	var contentType string
	if formdata {
		writer := multipart.NewWriter(&buf)
		err := network.parseFormData(body, writer, files)
		if err != nil {
			return nil, nil, err
		}
		err = writer.Close()
		if err != nil {
			logger.Error(errors.New("error closing writer after parsing formdata"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return nil, nil, err
		}
		contentType = writer.FormDataContentType()
	} else {
		parsedBody, err := json.Marshal(body)
		if err != nil {
			logger.Error(errors.New("error converting body to JSON"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return nil, nil, err
		}
		buf = *bytes.NewBuffer(parsedBody)
		contentType = "application/json"
	}

	req, err := http.NewRequest("POST", network.BaseUrl+path, &buf)
	req.Header.Set("Content-Type", contentType)
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

func (network *NetworkController) parseFormData(body any, writer *multipart.Writer, files *map[string][]byte) error {
	if body == nil {
		return errors.New("empty body passed")
	}

	for i, item := range any(body).(map[string]any) {
		field, err := writer.CreateFormField(i)
		if err != nil {
			logger.Error(errors.New("error parsing multipart payload"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return err
		}
		b, err := json.Marshal(item)
		if err != nil {
			logger.Error(errors.New("error marshalling item to byte for formdata"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
		}
		_, err = io.Copy(field, bytes.NewReader(b))
		if err != nil {
			logger.Error(errors.New("error copying file to file part for formdata"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return err
		}
	}
	for i, file := range *files {
		part, err := writer.CreateFormFile(i, utils.GenerateUUIDString())
		if err != nil {
			logger.Error(errors.New("error creating form file for formdata"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return err
		}
		if _, err := io.Copy(part, bytes.NewReader(file)); err != nil {
			logger.Error(errors.New("error copying file to file part for formdata"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			})
			return err
		}
	}
	return nil
}
