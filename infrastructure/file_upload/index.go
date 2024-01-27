package fileupload

import (
	"os"

	"kego.com/infrastructure/file_upload/azure"
	"kego.com/infrastructure/file_upload/types"
)

var FileUploader types.FileUploaderType

func InitialiseFileUploader(){
	FileUploader = &azure.AzureBlobSignedURLService{
		AccountName: os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AccountKey: os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		ContainerName: os.Getenv("AZURE_CONTAINER_NAME"),
	}
}