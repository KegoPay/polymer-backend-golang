package entities

import (
	"kego.com/application/utils"
)

type ExchangeRates struct{
	USDNGN 		 float32		`json:"USDNGN"`
	NGNUSD 		 float32
	USDCAD 		 float32		`json:"USDCAD"`
	NGNCAD 		 float32
	USDGHS 		 float32		`json:"USDGHS"`
	NGNGHS 		 float32
	USDINR 		 float32		`json:"USDINR"`
	NGNINR 		 float32
	USDKES 		 float32		`json:"USDKES"`
	NGNKES 		 float32
	USDZAR 		 float32		`json:"USDZAR"`
	NGNZAR 		 float32
	USDGBP 		 float32		`json:"USDGBP"`
	NGNGBP 		 float32

}

type ParsedExchangeRates struct {
	USDRate float32 `json:"USDrate"`
	NGNRate float32 `json:"NGNrate"`
}

func (er *ExchangeRates) FormatAllRates(intAmount *uint64) *map[string]ParsedExchangeRates {
	var amount float32
	if intAmount == nil  || *intAmount == 0 {
		amount = 1.0
	}else {
		amount = utils.UInt64ToFloat32Currency(*intAmount)
	}
	var formatedRates  map[string]ParsedExchangeRates = map[string]ParsedExchangeRates{}
	formatedRates["American Dollar (US) - Naira"] = ParsedExchangeRates{
		USDRate: amount,
		NGNRate: (er.USDNGN) * amount,
	}
	formatedRates["Canadian Dollar (CA) - Naira"] = ParsedExchangeRates{
		USDRate: (1/er.USDCAD * amount),
		NGNRate: (1/er.USDCAD * er.USDNGN) * amount,
	}
	formatedRates["British Pounds (GB) - Naira"] = ParsedExchangeRates{
		USDRate: (1/er.USDGBP * amount),
		NGNRate: (1/er.USDGBP * er.USDNGN) * amount,
	}
	formatedRates["South African Rand (ZA)- Naira"] = ParsedExchangeRates{
		USDRate: (1/er.USDZAR * amount),
		NGNRate: (1/er.USDZAR * er.USDNGN) * amount,
	}
	formatedRates["Ghanaian Cedis (GH) - Naira"] =  ParsedExchangeRates{
		USDRate: (1/er.USDGHS * amount),
		NGNRate: (1/er.USDGHS * er.USDNGN) * amount,
	}
	formatedRates["Indian Rupees (IN) - Naira"] =  ParsedExchangeRates{
		USDRate: (1/er.USDINR * amount),
		NGNRate: (1/er.USDINR * er.USDNGN) * amount,
	}
	formatedRates["Kenyan Shilling (KE) - Naira"] =  ParsedExchangeRates{
		USDRate: (1/er.USDKES * amount),
		NGNRate: (1/er.USDKES * er.USDNGN) * amount,
	}
	 return &formatedRates
}