package types

type FileUploaderType interface {
	GeneratedSignedURL(file_name string, permission SignedURLPermission) (*string, error)
	DeleteFileByURL(file_url string) (bool, error)
}

type SignedURLPermission struct {
	Read bool
	Write bool
}