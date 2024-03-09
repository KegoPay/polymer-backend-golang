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
		logger.Error(errors.New("error retireving bvn data from dojah"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for BVN fetch was unsuccessful"), logger.LoggerOptions{
			Key: "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key: "data",
			Data: dojahResponse,
		})
		return nil, errors.New("error retireving bvn")
	}
	logger.Info("NIN information retireved by Dojah")
	return &dojahResponse.Data, nil
}

func (div *DojahIdentityVerification) FetchNINDetails(nin string) (*identity_verification_types.NINData, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/kyc/nin/advance?nin=%s", nin), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId": div.APP_ID,
	}, nil)
	var dojahResponse DojahNINResponse
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from dojah"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for BVN fetch was unsuccessful"), logger.LoggerOptions{
			Key: "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key: "data",
			Data: dojahResponse,
		})
		return nil, errors.New("error retireving bvn")
	}
	logger.Info("BVN information retireved by Dojah")
	return &dojahResponse.Data, nil
}


func (div *DojahIdentityVerification) EmailVerification(email string) (bool, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/fraud/email?email_address=%s",email) , &map[string]string{
		"Authorization": div.API_KEY,
		"AppId": div.APP_ID,
	}, nil)
	var dojahResponse DojahEmailVerification
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error verifying email from dojah"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return false, errors.New("something went wrong while verifying email from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah email verification was unsuccessful"), logger.LoggerOptions{
			Key: "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key: "data",
			Data: dojahResponse,
		})
		return false, errors.New("error verifying email")
	}
	logger.Info("Email verification successful", logger.LoggerOptions{
		Key: "email",
		Data: email,
	}, logger.LoggerOptions{
		Key: "result",
		Data: dojahResponse,
	})
	return dojahResponse.Entity.Deliverable && !dojahResponse.Entity.DomainDetails.SusTLD && dojahResponse.Entity.DomainDetails.Registered && (dojahResponse.Entity.Score == 1), nil
}