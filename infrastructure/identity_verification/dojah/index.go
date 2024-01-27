package dojah_identity_verification

import (
	"encoding/json"
	"errors"
	"fmt"

	identity_verification_types "kego.com/infrastructure/identity_verification/types"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)

type DojahIdentityVerification struct {
	Network *network.NetworkController
	API_KEY string
	APP_ID string
}

func (div *DojahIdentityVerification) FetchBVNDetails(bvn string) (*identity_verification_types.BVNData, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/kyc/bvn/advance?bvn=%s", bvn), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId": div.APP_ID,
	}, nil)
	var dojahResponse DojahBVNResponse
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from prembly"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from prembly")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for BVN fetch was unsuccessful"), logger.LoggerOptions{
			Key: "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key: "data",
			Data: dojahResponse,
		})
		return nil, errors.New("error retrivin bvn")
	}
	logger.Info("BVN information retireved by Dojah")
	return &dojahResponse.Data, nil
}
