package international_payment_processor

import "kego.com/entities"


type ChimoneyExchangeRateDTO struct {
	Status 	bool						  `json:"status"`
	Data	entities.ExchangeRates		  `json:"data"`
}