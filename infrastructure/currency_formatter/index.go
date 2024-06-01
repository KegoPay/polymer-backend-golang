package currencyformatter

import (
	"github.com/dustin/go-humanize"
	"usepolymer.co/application/utils"
)

func HumanReadableFloat32Currency(amount float32) string {
	return humanize.FormatFloat("#,###.##", float64(amount))
}

func HumanReadableIntCurrency(amount uint64) string {
	return humanize.FormatFloat("#,###.##", float64(utils.UInt64ToFloat32Currency(amount)))
}
