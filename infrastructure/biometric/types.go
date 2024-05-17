package biometric

import faceapi "kego.com/infrastructure/biometric/faceAPI"

type BiometricServiceType interface {
	FaceMatchWithLiveness(referenceImg []byte, deviceID string) (*faceapi.FaceAPIFaceMatchWithLivenessResponse, error)
}