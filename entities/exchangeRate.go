package entities

import (
	"kego.com/application/utils"
)

type ExchangeRates struct{
	USDNGN 		 float32		`json:"USDNGN"`
	USDCAD 		 float32		`json:"USDCAD"`
	USDGHS 		 float32		`json:"USDGHS"`
	USDINR 		 float32		`json:"USDINR"`
	USDKES 		 float32		`json:"USDKES"`
	USDZAR 		 float32		`json:"USDZAR"`
	USDGBP 		 float32		`json:"USDGBP"`

}

type ParsedExchangeRates struct {
	USDRate float32 `json:"USDrate"`
	NGNRate float32 `json:"NGNrate"`
}

func (er *ExchangeRates) FormatAllRates(amount *uint64) map[string]ParsedExchangeRates {
	if amount == nil  || *amount == 0 {
		amount = utils.GetUInt64Pointer(1)
	}
	var formatedRates  map[string]ParsedExchangeRates = map[string]ParsedExchangeRates{}
	formatedRates["American Dollar (US) - Naira"] = ParsedExchangeRates{
		USDRate: 1 * float32(*amount),
		NGNRate: er.USDNGN * float32(*amount),
	}
	formatedRates["Canadian Dollar (CA) - Naira"] = ParsedExchangeRates{
		USDRate: er.USDCAD * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDCAD) * float32(*amount),
	}
	formatedRates["British Pounds (GB) - Naira"] = ParsedExchangeRates{
		USDRate: er.USDGBP * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDGBP) * float32(*amount),
	}
	formatedRates["South African Rand (ZA)- Naira"] = ParsedExchangeRates{
		USDRate: er.USDZAR * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDZAR) * float32(*amount),
	}
	formatedRates["Ghanaian Cedis (GH) - Naira"] =  ParsedExchangeRates{
		USDRate: er.USDGHS * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDGHS) * float32(*amount),
	}
	formatedRates["Indian Rupees (IN) - Naira"] =  ParsedExchangeRates{
		USDRate: er.USDINR * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDINR) * float32(*amount),
	}
	formatedRates["Kenyan Shilling (KE) - Naira"] =  ParsedExchangeRates{
		USDRate: er.USDKES * float32(*amount),
		NGNRate: er.FormatAgainstNGN(er.USDNGN, er.USDKES) * float32(*amount),
	}
	 return formatedRates
}

func (er *ExchangeRates) FormatAgainstNGN(NGNUSDRate float32, otherRate float32) float32 {
	return NGNUSDRate * (1 / otherRate)
}