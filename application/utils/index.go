package utils

import (
	"math"
	"regexp"

	"github.com/google/uuid"
)

func GenerateUUIDString() string {
	return uuid.NewString()
}

func ParseAmountToSmallerDenomination(amount uint64) uint64 {
	return amount * 100
}

func ParseAmountToHigherDenomination(amount uint64) uint64 {
	return amount / 100
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

func Float32ToUint64Currency(value float32) uint64 {
	roundedValue := math.Round(float64(value * 100))
	uintValue := uint64(roundedValue)
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
