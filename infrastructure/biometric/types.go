package biometric

import faceapi "usepolymer.co/infrastructure/biometric/faceAPI"

type BiometricServiceType interface {
	FaceMatchWithLiveness(referenceImg []byte, deviceID string) (*faceapi.FaceAPIFaceMatchWithLivenessResponse, error)
}
