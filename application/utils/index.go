package utils

import (
	"regexp"
	"time"

	"github.com/oklog/ulid/v2"
	"kego.com/application/constants"
)


func GenerateUUIDString() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulid.DefaultEntropy()).String()
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

func GetInt64Pointer(data int64) *int64 {
	return &data
}

func GetInternationalTransactionFee(amount float32) (internationalProcessorFee float32, transactionFee float32, transactionFeeVat float32) {
	transactionFee = amount * constants.INTERNATIONAL_TRANSACTION_FEE_RATE
	transactionFeeVat = transactionFee * constants.INTERNATIONAL_TRANSACTION_FEE_VAT
	internationalProcessorFee = amount * constants.INTERNATIONAL_PROCESSOR_FEE_RATE
    return
}

func GetLocalTransactionFee(amount uint64) (localProcessorFee float32, localProcessorVAT float32, polymerFee float32, polymerVAT float32) {
    if amount <= 500000 {
		polymerFee = constants.LOCAL_PROCESSOR_FEE_LT_5000 * constants.LOCAL_TRANSACTION_FEE_RATE
		polymerVAT = polymerFee * constants.LOCAL_TRANSACTION_FEE_VAT
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_LT_5000 
		localProcessorVAT = constants.LOCAL_PROCESSOR_FEE_LT_5000 * constants.LOCAL_TRANSACTION_FEE_VAT
    } else if amount <= 5000000 {
		polymerFee = constants.LOCAL_PROCESSOR_FEE_LT_50000 * constants.LOCAL_TRANSACTION_FEE_RATE
		polymerVAT = polymerFee * constants.LOCAL_TRANSACTION_FEE_VAT
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_LT_50000
		localProcessorVAT = constants.LOCAL_PROCESSOR_FEE_LT_50000 * constants.LOCAL_TRANSACTION_FEE_VAT
    } else {
		polymerFee = constants.LOCAL_PROCESSOR_FEE_GT_50000 * constants.LOCAL_TRANSACTION_FEE_RATE
		polymerVAT = polymerFee * constants.LOCAL_TRANSACTION_FEE_VAT
        localProcessorFee = constants.LOCAL_PROCESSOR_FEE_GT_50000
		localProcessorVAT = constants.LOCAL_PROCESSOR_FEE_GT_50000 * constants.LOCAL_TRANSACTION_FEE_VAT
    }
    return
}

func CountryCodeToCountryName(code string) string {
	countryCodeMap := map[string]string {
		"NG": "Nigeria",
		"IN": "India",
		"US": "United States of America",
		"KE": "Kenya",
		"ZA": "South Africa",
		"GB": "Britain",
		"GH": "Ghana",
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

func CurrencyCodeToCountryCode(code string) string {
	countryCodeMap := map[string]string {
		"NGN": "NG",
		"INR": "IN",
		"USD": "US",
		"KES": "KE",
		"ZAR": "ZA",
		"GBP": "GB",
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
        "GHS": "GH₵",
    }
    return currencySymbolMap[code]
}


func Float32ToUint64Currency(value float32) uint64 {
	roundUp := (uint64(value * 10000000) % 100000) != 0
	uintValue := uint64(value * 100)
	if roundUp {
		uintValue++
	}
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
