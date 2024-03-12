package entities

import (
	"kego.com/application/utils"
)

type ExchangeRates struct{
	USDNGN 		 float32		`json:"USDNGN"`
	NGNUSD 		 float32		`json:"NGNUSD"`
	USDCAD 		 float32		`json:"USDCAD"`
	NGNCAD 		 float32		`json:"NGNCAD"`
	USDGHS 		 float32		`json:"USDGHS"`
	NGNGHS 		 float32		`json:"NGNGHS"`
	USDINR 		 float32		`json:"USDINR"`
	NGNINR 		 float32		`json:"NGNINR"`
	USDKES 		 float32		`json:"USDKES"`
	NGNKES 		 float32		`json:"NGNKES"`
	USDZAR 		 float32		`json:"USDZAR"`
	NGNZAR 		 float32		`json:"NGNZAR"`
	USDGBP 		 float32		`json:"USDGBP"`
	NGNGBP 		 float32		`json:"NGNGBP"`

}

type ParsedExchangeRates struct {
	USDRate float32 `json:"USDrate"`
	NGNRate float32 `json:"NGNrate"`
}

func (er *ExchangeRates) FormatAllRates(amount *uint64) *map[string]ParsedExchangeRates {
	if amount == nil  || *amount == 0 {
		amount = utils.GetUInt64Pointer(1)
	}
	var formatedRates  map[string]ParsedExchangeRates = map[string]ParsedExchangeRates{}
	formatedRates["American Dollar (US) - Naira"] = ParsedExchangeRates{
		USDRate: 1 * float32(*amount),
		NGNRate: (1/er.NGNUSD) * float32(*amount),
	}
	formatedRates["Canadian Dollar (CA) - Naira"] = ParsedExchangeRates{
		USDRate: (er.USDCAD * float32(*amount)),
		NGNRate: (1/er.NGNCAD) * float32(*amount),
	}
	formatedRates["British Pounds (GB) - Naira"] = ParsedExchangeRates{
		USDRate: (er.USDGBP * float32(*amount)),
		NGNRate: (1/er.NGNGBP) * float32(*amount),
	}
	formatedRates["South African Rand (ZA)- Naira"] = ParsedExchangeRates{
		USDRate: (er.USDZAR * float32(*amount)),
		NGNRate: (1/er.NGNZAR) * float32(*amount),
	}
	formatedRates["Ghanaian Cedis (GH) - Naira"] =  ParsedExchangeRates{
		USDRate: (er.USDGHS * float32(*amount)),
		NGNRate: (1/er.NGNGHS) * float32(*amount),
	}
	formatedRates["Indian Rupees (IN) - Naira"] =  ParsedExchangeRates{
		USDRate: (er.USDINR * float32(*amount)),
		NGNRate: (1/er.NGNINR) * float32(*amount),
	}
	formatedRates["Kenyan Shilling (KE) - Naira"] =  ParsedExchangeRates{
		USDRate: (er.USDKES * float32(*amount)),
		NGNRate: (1/er.NGNKES) * float32(*amount),
	}
	 return &formatedRates
}