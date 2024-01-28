package prembly_identity_verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"kego.com/application/constants"
	"kego.com/application/utils"
	identity_verification_types "kego.com/infrastructure/identity_verification/types"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)

type PremblyIdentityVerification struct {
	Network *network.NetworkController
	API_KEY string
	APP_ID string
}

func (piv *PremblyIdentityVerification) FetchBVNDetails(bvn string) (*identity_verification_types.BVNData, error) {
	response, _, err := piv.Network.Post("/identitypass/verification/bvn", &map[string]string{
		"x-api-key": piv.API_KEY,
		"app-id": piv.APP_ID,
	}, map[string]string{
		"number": bvn,
	}, nil)
	var premblyResponse PremblyBVNResponse
	json.Unmarshal(*response, &premblyResponse)
	fmt.Println(premblyResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from prembly"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from prembly")
	}
	if !premblyResponse.Status {
		logger.Error(errors.New(premblyResponse.Message), logger.LoggerOptions{
			Key: "error",
			Data: errors.New(premblyResponse.Detail),
		}, logger.LoggerOptions{
			Key: "data",
			Data: premblyResponse,
		})
		return nil, errors.New(premblyResponse.Message)
	}
	logger.Info("BVN information retireved by Prembly")
	return &premblyResponse.Data, nil
}

func (piv *PremblyIdentityVerification) FaceMatch(img1 string, img2 string) (*float32, error) {
	if os.Getenv("GIN_MODE") != "production" {
		img1 = "https://res.cloudinary.com/dh3i1wodq/image/upload/v1675417496/cbimage_3_drqdoc.jpg"
		img2 = "https://res.cloudinary.com/dh3i1wodq/image/upload/v1677955197/face_image_tkmmwz.jpg"
	}
	response, _, err := piv.Network.Post("/identitypass/verification/biometrics/face/comparison", &map[string]string{
		"x-api-key": piv.API_KEY,
		"app-id": piv.APP_ID,
	}, map[string]string{
		"image_one": img1,
		"image_two": img2,
	}, nil)
	var premblyResponse PremblyFaceMatchResponse
	json.Unmarshal(*response, &premblyResponse)
	if err != nil {
		logger.Error(errors.New("error performing face match on prembly"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while performing face match")
	}
	if !premblyResponse.Status {
		logger.Error(errors.New(premblyResponse.Message), logger.LoggerOptions{
			Key: "error",
			Data: errors.New(premblyResponse.Detail),
		})
		if (premblyResponse.Detail == "Face does not match") {
			premblyResponse.Detail = fmt.Sprintf("Your picture does not match with the Image on the BVN provided. If you think this is a mistake please contact support on %s", constants.SUPPORT_EMAIL)
		}
		return nil, errors.New(premblyResponse.Detail)
	}
	if premblyResponse.Detail != "" {
		return utils.GetFloat32Pointer(0), nil
	}
	logger.Info("Face Match completed by Prembly")
	return &premblyResponse.Confidence, nil
}

func (piv *PremblyIdentityVerification) ValidateEmail(bvn string) (*identity_verification_types.BVNData, error) {
	response, _, err := piv.Network.Post("/identityradar/api/v1/email-intelligence", &map[string]string{
		"api-key": piv.API_KEY,
	}, map[string]string{
		"number": bvn,
	}, nil)
	var premblyResponse PremblyBVNResponse
	json.Unmarshal(*response, &premblyResponse)
	fmt.Println(premblyResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from prembly"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from prembly")
	}
	if !premblyResponse.Status {
		logger.Error(errors.New(premblyResponse.Message), logger.LoggerOptions{
			Key: "error",
			Data: errors.New(premblyResponse.Detail),
		}, logger.LoggerOptions{
			Key: "data",
			Data: premblyResponse,
		})
		return nil, errors.New(premblyResponse.Message)
	}
	logger.Info("BVN information retireved by Prembly")
	return &premblyResponse.Data, nil
}