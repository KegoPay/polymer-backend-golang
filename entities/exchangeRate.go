package entities

import (
	"fmt"
	"strconv"
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

func (er *ExchangeRates) FormatAllRates() map[string]float32 {
	formatedRates := map[string]float32{}
	formatedRates["American Dollar (US) - Naira"] = er.USDNGN
	formatedRates["Canadian Dollar (CA) - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDCAD)
	formatedRates["British Pounds - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDGBP)
	formatedRates["South African Rand (ZA)- Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDZAR)
	formatedRates["Ghanaian Cedis (GH) - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDGHS)
	formatedRates["Indian Rupees (IN) - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDINR)
	formatedRates["Kenyan Shilling (KE) - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDKES)
	 return formatedRates
}

func (er *ExchangeRates) FormatAgainstNGN(ngnusdRate float32, usdotherRate float32) float32 {
	rate := ngnusdRate * (1 / usdotherRate)
	stringRate := fmt.Sprintf("%.2f", rate)
	parsedRate, _ := strconv.ParseFloat(stringRate, 32)
	return float32(parsedRate)
}