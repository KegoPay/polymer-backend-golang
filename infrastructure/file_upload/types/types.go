package types

type FileUploaderType interface {
	GeneratedSignedURL(file_name string, permission SignedURLPermission) (*string, error)
	UploadBase64File(file_name string, file *string) error
}

type SignedURLPermission struct {
	Read bool `json:"read"`
	Write bool `json:"write"`
}