package fileupload

import (
	"errors"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	cloudinaryService "kego.com/infrastructure/file_upload/cloudinary"
	"kego.com/infrastructure/file_upload/types"
	"kego.com/infrastructure/logger"
)

var FileUploader types.FileUploaderType

func InitialiseFileUploader(){
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	if err != nil {
		logger.Error(errors.New("error creating instance of cloudinary"), logger.LoggerOptions{
			Key: "err",
			Data: err,
		})
	}
	FileUploader = &cloudinaryService.CloudinaryFileUpload{
		CloudinaryConfig: cld,
	}
}