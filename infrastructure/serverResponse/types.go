package server_response

type serverResponder interface{
	// Used to send a JSON response to the client.
	Respond(ctx interface{}, code int, message string, payload interface{}, errs []error, response_code *uint, device_id *string)
	UnEncryptedRespond(ctx interface{}, code int, message string, payload interface{}, errs []error, response_code *uint)
}