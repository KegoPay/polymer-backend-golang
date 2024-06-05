package dojah_identity_verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	identity_verification_types "usepolymer.co/infrastructure/identity_verification/types"
	"usepolymer.co/infrastructure/logger"
	"usepolymer.co/infrastructure/network"
)

type DojahIdentityVerification struct {
	Network *network.NetworkController
	API_KEY string
	APP_ID  string
}

func (div *DojahIdentityVerification) FetchAdvancedCACDetails(rc string) (*identity_verification_types.CompanyProfile, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/kyc/cac/advance?class=premium&type=bn&rc=%s", rc), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId":         div.APP_ID,
	}, nil)
	if err != nil {
		logger.Error(errors.New("error retireving cac data from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving cac data from dojah")
	}
	var dojahResponse identity_verification_types.CACResponse
	err = json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error parsing cac data from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while parsing cac data from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for cac fetch was unsuccessful"), logger.LoggerOptions{
			Key:  "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: dojahResponse,
		})
		return nil, errors.New("error retireving cac")
	}
	logger.Info("cac information retireved by Dojah")
	if os.Getenv("ENV") != "prod" {
		var mockResponse identity_verification_types.CompanyProfile
		json.Unmarshal(PolymerCACDetails, &mockResponse)
		dojahResponse.Data = mockResponse
	}
	return &dojahResponse.Data, nil
}

func (div *DojahIdentityVerification) ImgLivenessCheck(img string) (bool, error) {
	response, statusCode, err := div.Network.Post("/ml/liveness/", &map[string]string{
		"Authorization": div.API_KEY,
		"AppId":         div.APP_ID,
	}, map[string]any{
		"image": img,
	}, nil, false, nil)
	if err != nil {
		logger.Error(errors.New("error liveness check result from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, errors.New("something went wrong while retireving liveness check result from dojah")
	}
	var dojahResponse identity_verification_types.LivenessCheckResult
	err = json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error parsing liveness check result from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, errors.New("something went wrong while parsing liveness check result from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for liveness check was unsuccessful"), logger.LoggerOptions{
			Key:  "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: dojahResponse,
		})
		return false, errors.New("error retireving liveness check result ")
	}
	logger.Info("liveness check completed by Dojah")
	if !dojahResponse.Entity.Liveness.LivenessCheck {
		return false, errors.New("Face verification failed. Please ensure you are in a well lit environment and have no coverings on your face.")
	}
	if dojahResponse.Entity.Liveness.LivenessProbability < 60.0 {
		return false, errors.New("Face verification failed. Please ensure you are in a well lit environment and have no coverings on your face.")
	}
	return true, nil
}

func (div *DojahIdentityVerification) FetchBVNDetails(bvn string) (*identity_verification_types.BVNData, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/kyc/bvn/advance?bvn=%s", bvn), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId":         div.APP_ID,
	}, nil)
	var dojahResponse DojahBVNResponse
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving bvn data from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for BVN fetch was unsuccessful"), logger.LoggerOptions{
			Key:  "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: dojahResponse,
		})
		return nil, errors.New("error retireving bvn")
	}
	logger.Info("BVN information retireved by Dojah")
	return &dojahResponse.Data, nil
}

func (div *DojahIdentityVerification) FetchNINDetails(nin string) (*identity_verification_types.NINData, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/kyc/nin/advance?nin=%s", nin), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId":         div.APP_ID,
	}, nil)
	var dojahResponse DojahNINResponse
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error retireving nin data from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, errors.New("something went wrong while retireving nin data from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for nin fetch was unsuccessful"), logger.LoggerOptions{
			Key:  "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: dojahResponse,
		})
		if dojahResponse.Error == "Wrong NIN Inputted" {
			return nil, errors.New("NIN not found. Crosscheck the number inputed")
		}
		return nil, errors.New("error retireving nin")
	}
	logger.Info("NIN information retireved by Dojah")
	return &dojahResponse.Data, nil
}

func (div *DojahIdentityVerification) EmailVerification(email string) (bool, error) {
	response, statusCode, err := div.Network.Get(fmt.Sprintf("/fraud/email?email_address=%s", email), &map[string]string{
		"Authorization": div.API_KEY,
		"AppId":         div.APP_ID,
	}, nil)
	var dojahResponse DojahEmailVerification
	json.Unmarshal(*response, &dojahResponse)
	if err != nil {
		logger.Error(errors.New("error verifying email from dojah"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return false, errors.New("something went wrong while verifying email from dojah")
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah email verification was unsuccessful"), logger.LoggerOptions{
			Key:  "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: dojahResponse,
		})
		return false, errors.New("error verifying email")
	}
	logger.Info("Email verification successful", logger.LoggerOptions{
		Key:  "email",
		Data: email,
	}, logger.LoggerOptions{
		Key:  "result",
		Data: dojahResponse,
	})
	return dojahResponse.Entity.Deliverable && !dojahResponse.Entity.DomainDetails.SusTLD && dojahResponse.Entity.DomainDetails.Registered && (dojahResponse.Entity.Score == 1), nil
}
