package fileupload

import (
	"os"

	"usepolymer.co/infrastructure/file_upload/azure"
	"usepolymer.co/infrastructure/file_upload/types"
)

var FileUploader types.FileUploaderType

func InitialiseFileUploader() {
	FileUploader = &azure.AzureBlobSignedURLService{
		AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AccountKey:    os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		ContainerName: os.Getenv("AZURE_CONTAINER_NAME"),
	}
}
