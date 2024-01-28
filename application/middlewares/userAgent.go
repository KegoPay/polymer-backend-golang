package middlewares

import (
	"errors"
	"os"
	"strings"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/ipresolver"
	"kego.com/infrastructure/logger"
	useragent "kego.com/infrastructure/user_agent"
)


func UserAgentMiddleware(ctx *interfaces.ApplicationContext[any], minAppVersion string, ipAddress string, mobileOnly bool) (*interfaces.ApplicationContext[any], bool) {
	agent := ctx.GetHeader("User-Agent")
	if agent == nil {
		apperrors.ClientError(ctx.Ctx, "Why your User-Agent header no dey? You be criminal?🤨", []error{errors.New("user agent header missing")})
		return nil, false
	}

	userAgentData := useragent.Parse(agent.(string))
	if os.Getenv("GIN_MODE") == "release" {
		if !userAgentData.IsMobile && !userAgentData.IsTablet && !userAgentData.IsDesktop {
			apperrors.UnsupportedUserAgent(ctx.Ctx)
			return nil ,false
		}
	}
	

	if mobileOnly && userAgentData.AgentName != "Polymer" {
		apperrors.UnsupportedUserAgent(ctx.Ctx)
		return nil ,false
	}

	if (mobileOnly && userAgentData.OSName != "Android" && userAgentData.OSName != "iOS"){
		apperrors.UnsupportedUserAgent(ctx.Ctx)
		return nil ,false
	}

	reqSemVers  := strings.Split(userAgentData.BuildNumber, ".")
	minAppVersionSemVers := strings.Split(minAppVersion, ".")
	if minAppVersionSemVers[0] > reqSemVers[0] {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return nil, false
	}
	if minAppVersionSemVers[1] > reqSemVers[1] {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return nil, false
	}
	if minAppVersionSemVers[2] > reqSemVers[2] {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return nil, false
	}
	ipLookupRes, err  := ipresolver.IPResolverInstance.LookUp(ipAddress)
	if err != nil {
		logger.Error(errors.New("error looking up ip"), logger.LoggerOptions{
			Key: "ip",
			Data: ipAddress,
		}, logger.LoggerOptions{
			Key: "user agent",
			Data: agent,
		})
	}
	logger.Info("request-ip-lookup", logger.LoggerOptions{
		Key: "ip-data",
		Data: ipLookupRes,
	}, logger.LoggerOptions{
		Key: "user-agent",
		Data: agent,
	})

	return ctx, true
}