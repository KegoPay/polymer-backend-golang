package entities

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
	formatedRates["United States Dollar - Naira"] = er.USDNGN
	formatedRates["Canadian Dollar - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDCAD)
	formatedRates["British Pounds - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDGBP)
	formatedRates["South African Rand - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDZAR)
	formatedRates["Ghana Cedis - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDGHS)
	formatedRates["Indian Rupees - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDINR)
	formatedRates["Kenyan Shilling - Naira"] = er.FormatAgainstNGN(er.USDNGN, er.USDKES)
	 return formatedRates
}

func (er *ExchangeRates) FormatAgainstNGN(ngnusdRate float32, usdotherRate float32) float32 {
	return ngnusdRate * (1 / usdotherRate)
}