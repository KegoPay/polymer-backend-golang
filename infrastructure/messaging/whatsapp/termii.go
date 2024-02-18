package sms

import (
	"encoding/json"
	"errors"
	"fmt"

	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/network"
)

type TermiiService struct {
	Network *network.NetworkController
	API_KEY string
}

func (ts *TermiiService) SendSMS(phone string, sms string) *string {
	response, statusCode, err := ts.Network.Post("/sms/otp/send", nil, map[string]any{
		"api_key": ts.API_KEY,
		"message_type": "NUMERIC",
		"from": "Polymer",
		"to": phone,
		"channel": "generic", // temp
		"pin_attempts": 3,
		"pin_time_to_live": 10,
		"pin_length": 6,
		"pin_placeholder": "< 123456 >",
		"message_text": sms,
	}, nil)
	fmt.Println(err)
	fmt.Println(phone)
	fmt.Println(*statusCode)
	var termiiResponse TermiiOTPResponse
	json.Unmarshal(*response, &termiiResponse)
	fmt.Println(termiiResponse)
	if err != nil {
		logger.Error(errors.New("error retireving bvn data from dojah"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	if *statusCode != 200 {
		logger.Error(errors.New("request to Dojah for BVN fetch was unsuccessful"), logger.LoggerOptions{
			Key: "statusCode",
			Data: fmt.Sprintf("%d", statusCode),
		}, logger.LoggerOptions{
			Key: "data",
			Data: termiiResponse,
		})
		return nil
	}
	logger.Info(fmt.Sprintf("SMS OTP sent to %s", phone), logger.LoggerOptions{
		Key: "res",
		Data: termiiResponse,
	})
	return &termiiResponse.PinID
}