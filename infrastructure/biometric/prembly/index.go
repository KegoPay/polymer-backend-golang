package prembly_idpass

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"usepolymer.co/application/utils"
	"usepolymer.co/infrastructure/biometric/types"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

type PremblyBiometricService struct {
	Network *network.NetworkController
	API_KEY string
	APP_ID  string
}

// can either be a url or base64 string
func (piv *PremblyBiometricService) LivenessCheck(img *string) (bool, error) {
	if os.Getenv("ENV") != "prod" {
		img = utils.GetStringPointer("https://res.cloudinary.com/dh3i1wodq/image/upload/v1675417496/cbimage_3_drqdoc.jpg")
	}
	response, _, err := piv.Network.Post("/identitypass/verification/biometrics/face/liveliness_check", &map[string]string{
		"x-api-key": piv.API_KEY,
		"app-id":    piv.APP_ID,
	}, map[string]any{
		"image": img,
	}, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("error performing liveness check on Prembly"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, err
	}
	var premblyResponse types.BiometricLivenessResponse
	err = json.Unmarshal(*response, &premblyResponse)
	if err != nil {
		err = errors.New("error parsing response for liveness check response on Prembly")
		logger.Error(err, logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, err
	}
	logger.Info("Liveness check completed by Prembly", logger.LoggerOptions{
		Key:  "message",
		Data: premblyResponse.Message,
	})
	return premblyResponse.Result > 75.0, nil
}

// can either be a url or base64 string
func (piv *PremblyBiometricService) FaceMatch(img1 *string, img2 *string) (bool, error) {
	if os.Getenv("ENV") != "prod" {
		img1 = utils.GetStringPointer("https://res.cloudinary.com/dh3i1wodq/image/upload/v1675417496/cbimage_3_drqdoc.jpg")
		img2 = utils.GetStringPointer("https://res.cloudinary.com/dh3i1wodq/image/upload/v1677955197/face_image_tkmmwz.jpg")
	}
	response, _, err := piv.Network.Post("/identitypass/verification/biometrics/face/comparison", &map[string]string{
		"x-api-key": piv.API_KEY,
		"app-id":    piv.APP_ID,
	}, map[string]any{
		"image_one": img1,
		"image_two": img2,
	}, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("error performing face match check on Prembly"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, err
	}
	var premblyResponse types.BiometricFaceMatchResponse
	var premblyResponsee any
	err = json.Unmarshal(*response, &premblyResponse)
	err = json.Unmarshal(*response, &premblyResponsee)
	fmt.Println(premblyResponsee)
	if err != nil {
		logger.Error(errors.New("error parsing response for face match check response on Prembly"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, errors.New("error parsing response for face match check response on Prembly")
	}
	logger.Info("Face match check completed by Prembly", logger.LoggerOptions{
		Key:  "message",
		Data: premblyResponse.Message,
	})
	return premblyResponse.Score > 75.0, nil
}
