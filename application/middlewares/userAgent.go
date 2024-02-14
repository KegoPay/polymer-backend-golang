package middlewares

import (
	"errors"
	"regexp"
	"strings"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/ipresolver"
	"kego.com/infrastructure/logger"
)


func UserAgentMiddleware(ctx *interfaces.ApplicationContext[any], minAppVersion string, clientIP string) (*interfaces.ApplicationContext[any], bool) {
	agent := ctx.GetHeader("User-Agent")
	if agent == nil {
		apperrors.ClientError(ctx.Ctx, "Why your User-Agent header no dey? You be criminal?🤨", []error{errors.New("user agent header missing")})
		return nil, false
	}

	if !strings.Contains(agent.(string), "Android") && !strings.Contains(agent.(string), "iOS") {
		apperrors.UnsupportedUserAgent(ctx.Ctx)
		return nil ,false
	}

	if !strings.Contains(agent.(string), "Polymer/") {
		apperrors.UnsupportedUserAgent(ctx.Ctx)
		return nil ,false
	}

	versionRegex := regexp.MustCompile(`Polymer/([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := versionRegex.FindStringSubmatch(agent.(string))

	if len(matches) != 2 {
		apperrors.UnsupportedAppVersion(ctx.Ctx)
		return nil, false
	}

	appVersion := matches[1]
	reqSemVers  := strings.Split(appVersion, ".")
	if len(reqSemVers) < 3 {

	}
	minAppVersionSemVers := strings.Split(minAppVersion, ".")
	if len(minAppVersionSemVers) < 3 {

	}
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
	
	ipLookupRes, err  := ipresolver.IPResolverInstance.LookUp(clientIP)
	if err != nil {
		logger.Error(errors.New("error looking up ip"), logger.LoggerOptions{
			Key: "ip",
			Data: clientIP,
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