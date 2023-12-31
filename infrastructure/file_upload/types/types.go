package types

import "mime/multipart"

type FileUploaderType interface {
	UploadSingleFile(*multipart.FileHeader, *string) (*string, error)
	DeleteSingleFile(string) (error)
}
