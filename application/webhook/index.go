package webhook

import "usepolymer.co/application/interfaces"

func PaystackWebhook(ctx *interfaces.ApplicationContext[CustomerVerificationDTO]) {
	if ctx.Body.Event == "customeridentification.success" {

	} else if ctx.Body.Event == "customeridentification.failed" {

	}
}
