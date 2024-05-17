package services

import (
	"errors"
	"os"

	"kego.com/infrastructure/network"
)

var BackgroundServiceInstance *BackgroundService

type BackgroundService struct {
	Network	*network.NetworkController
	Apikey string
}

func InitialiseBackgroundService() {
	BackgroundServiceInstance =  &BackgroundService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("POLYMER_BACKGROUND_SERVICE"),
		},
		Apikey: os.Getenv("POLYMER_BACKGROUND_SERVICE_API_KEY"),
	}
}

func (bs *BackgroundService)RequestAccountStatementGeneration(walletID string, email string, start string, end string) error {
	_ , statusCode, err := bs.Network.Post("/wallet/request-statement", &map[string]string{
		"Api-Key": bs.Apikey,
	}, map[string]any {
		"walletID": walletID,
		"email": email,
		"start": start,
		"end": end,
	}, nil, false, nil)
	if err != nil {
		return err
	}
	if *statusCode != 200 {
		return errors.New("an unknown error occured while generating account statement")
	}
	return nil
}