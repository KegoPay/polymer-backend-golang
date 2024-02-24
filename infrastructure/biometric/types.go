package biometric

type BiometricServiceType interface {
	FaceMatch(img1 string, img2 string) (*float32, error)
}