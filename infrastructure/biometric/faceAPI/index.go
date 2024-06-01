package faceapi

import (
	"encoding/json"
	"errors"

	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

type FaceAPIBiometricService struct {
	Network *network.NetworkController
	API_KEY string
}

func (face *FaceAPIBiometricService) FaceMatchWithLiveness(referenceImg []byte, deviceID string) (*FaceAPIFaceMatchWithLivenessResponse, error) {
	response, _, err := face.Network.Post("/face/v1.1-preview.1/detectlivenesswithverify/singlemodal/sessions", &map[string]string{
		"Ocp-Apim-Subscription-Key": face.API_KEY,
	}, map[string]any{
		"Parameters": map[string]any{
			"livenessOperationMode": "passive",
			"deviceCorrelationId":   deviceID,
		},
	}, nil, true, &map[string][]byte{
		"VerifyImage": referenceImg,
	})
	if err != nil {
		logger.Error(errors.New("error creating FaceMatchWithLiveness session"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}
	var faceAPIResponse FaceAPIFaceMatchWithLivenessResponse
	err = json.Unmarshal(*response, &faceAPIResponse)
	if err != nil {
		err = errors.New("error parsing response for FaceMatchWithLiveness")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}
	logger.Info("Face Match session created by FaceAPI")
	return &faceAPIResponse, nil
}
