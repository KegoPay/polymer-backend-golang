package middlewares

import (
	"os"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/interfaces"
	"kego.com/infrastructure/cryptography"
)


func WebAgentMiddleware(ctx *interfaces.ApplicationContext[any], ipAddress string) (*interfaces.ApplicationContext[any], bool) {
	// agent := ctx.GetHeader("User-Agent")
	// if agent == nil {
	// 	apperrors.ClientError(ctx.Ctx, "Why your User-Agent header no dey? You be criminal?🤨", []error{errors.New("user agent header missing")})
	// 	return nil, false
	// }

	// if !strings.Contains(agent.(string), "Android") && !strings.Contains(agent.(string), "iOS") {
	// 	apperrors.UnsupportedUserAgent(ctx.Ctx)
	// 	return nil ,false
	// }

	// if !strings.Contains(agent.(string), "Polymer/") {
	// 	apperrors.UnsupportedUserAgent(ctx.Ctx)
	// 	return nil ,false
	// }

	// versionRegex := regexp.MustCompile(`Polymer/([0-9]+\.[0-9]+\.[0-9]+)`)
	// matches := versionRegex.FindStringSubmatch(agent.(string))

	// if len(matches) != 2 {
	// 	apperrors.UnsupportedAppVersion(ctx.Ctx)
	// 	return nil, false
	// }

	// appVersion := matches[1]
	// reqSemVers  := strings.Split(appVersion, ".")
	// if len(reqSemVers) < 3 {

	// }
	// minAppVersionSemVers := strings.Split(minAppVersion, ".")
	// if len(minAppVersionSemVers) < 3 {

	// }
	// if minAppVersionSemVers[0] > reqSemVers[0] {
	// 	apperrors.UnsupportedAppVersion(ctx.Ctx)
	// 	return nil, false
	// }
	// if minAppVersionSemVers[1] > reqSemVers[1] {
	// 	apperrors.UnsupportedAppVersion(ctx.Ctx)
	// 	return nil, false
	// }
	// if minAppVersionSemVers[2] > reqSemVers[2] {
	// 	apperrors.UnsupportedAppVersion(ctx.Ctx)
	// 	return nil, false
	// }
	// fmt.Println("here p")
	// ipLookupRes, err  := ipresolver.IPResolverInstance.LookUp("102.89.32.54")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(ipLookupRes)
	clientKey := ctx.GetHeader("Web_api_key")
	if clientKey == nil {
		apperrors.AuthenticationError(ctx.Ctx, "missing web token")
		return nil, false
	}
	valid := cryptography.CryptoHahser.VerifyData(os.Getenv("WEB_CLIENT_API_KEY"), clientKey.(string))
	if !valid {
		apperrors.AuthenticationError(ctx.Ctx, "invalid web token")
		return nil, false
	}
	return ctx, true
}