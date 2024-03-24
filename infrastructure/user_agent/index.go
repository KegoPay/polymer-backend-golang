package useragent

import (
	"github.com/mileusna/useragent"
)


func Parse(userAgent string) UserAgentData {
	parsed := useragent.Parse(userAgent)
	return UserAgentData {
		OSName: parsed.OS,
		AgentName: parsed.Name,
		BuildNumber: parsed.VersionNoFull(),
		DeviceName: parsed.Device,
		IsDesktop: parsed.Desktop,
		IsMobile: parsed.Mobile,
		IsTablet: parsed.Tablet,
	}
}