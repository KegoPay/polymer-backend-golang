package types

type BiometricServiceType interface {
	LivenessCheck(img *string) (bool, error)
	FaceMatch(img1 *string, img2 *string) (bool, error)
}

type BiometricLivenessResponse struct {
	Result  float32 `json:"confidence_in_percentage"`
	Message string  `json:"detail"`
}

type BiometricFaceMatchResponse struct {
	Score   float32   `json:"confidence"`
	Message string `json:"message"`
}
