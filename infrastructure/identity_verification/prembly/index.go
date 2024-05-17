package prembly_identity_verification

import (
	"encoding/json"
	"errors"

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
	}, map[string]any{
		"number": bvn,
	}, nil, false, nil)
	var premblyResponse PremblyBVNResponse
	json.Unmarshal(*response, &premblyResponse)
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


func (piv *PremblyIdentityVerification) ValidateEmail(bvn string) (*identity_verification_types.BVNData, error) {
	response, _, err := piv.Network.Post("/identityradar/api/v1/email-intelligence", &map[string]string{
		"api-key": piv.API_KEY,
	}, map[string]any{
		"number": bvn,
	}, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from prembly"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from prembly")
	}
	var premblyResponse PremblyBVNResponse
	err = json.Unmarshal(*response, &premblyResponse)
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