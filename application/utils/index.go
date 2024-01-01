package utils

import (
	"regexp"

	"github.com/google/uuid"
	"kego.com/application/constants"
)

func GenerateUUIDString() string {
	return uuid.NewString()
}

func GetStringPointer(text string) *string {
	return &text
}

func GetBooleanPointer(data bool) *bool {
	return &data
}

func GetFloat32Pointer(data float32) *float32 {
	return &data
}

func GetUInt64Pointer(data uint64) *uint64 {
	return &data
}

func GetTransactionFee(amount uint64) uint64 {
	fee := (float32(amount) * constants.INTERNATIONAL_TRANSACTION_FEE_RATE) / 100.0
    return Float32ToUint64Currency(fee)
}

func GetLocalTransactionFee(amount uint64) (localProcessorFee float32, polymerFee float32) {
	var vat float32
    if amount <= 500000 {
		vat = float32(constants.LOCAL_PROCESSOR_FEE_LT_5000) * constants.LOCAL_TRANSACTION_FEE_VAT
		polymerFee = float32(constants.LOCAL_PROCESSOR_FEE_LT_5000) * constants.LOCAL_TRANSACTION_FEE_RATE
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_LT_5000 + vat
    } else if amount <= 5000000 {
		vat = float32(constants.LOCAL_PROCESSOR_FEE_LT_50000) * constants.LOCAL_TRANSACTION_FEE_VAT
		polymerFee = float32(constants.LOCAL_PROCESSOR_FEE_LT_50000) * constants.LOCAL_TRANSACTION_FEE_RATE
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_LT_50000 + vat
    } else {
		vat = float32(constants.LOCAL_PROCESSOR_FEE_GT_50000) * constants.LOCAL_TRANSACTION_FEE_VAT
		polymerFee = float32(constants.LOCAL_PROCESSOR_FEE_GT_50000) * constants.LOCAL_TRANSACTION_FEE_RATE
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_GT_50000 + vat
    }
    return localProcessorFee, polymerFee
}

func CountryCodeToCountryName(code string) string {
	countryCodeMap := map[string]string {
		"NG": "Nigeria",
		"IN": "India",
		"US": "United States of America",
		"KE": "Kenya",
		"ZA": "South Africa",
		"GB": "Britain",
	}
	return countryCodeMap[code]
}

func CountryCodeToCurrencyCode(code string) string {
	countryCodeMap := map[string]string {
		"NG": "NGN",
		"IN": "INR",
		"US": "USD",
		"KE": "KES",
		"ZA": "ZAR",
		"GB": "GBP",
	}
	return countryCodeMap[code]
}

func CurrencyCodeToCurrencySymbol(code string) string {
    currencySymbolMap := map[string]string{
        "NGN": "₦",
        "INR": "₹",
        "USD": "$",
        "KES": "KSh",
        "ZAR": "R",
        "GBP": "£",
    }
    return currencySymbolMap[code]
}


func Float32ToUint64Currency(value float32) uint64 {
	uintValue := uint64(value * 100)
	return uintValue
}

func UInt64ToFloat32Currency(value uint64) float32 {
	floatValue := float32(value) / 100
	return floatValue
}

func ExtractAppVersionFromUserAgentHeader(userAgent string) *string {
	regex := regexp.MustCompile(`Polymer/([0-9.]+)`)
	matches := regex.FindStringSubmatch(userAgent)
	if len(matches) >= 2 {
		return &matches[1]
	}
	return nil
}
