package cloudinary

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"kego.com/application/utils"
	"kego.com/infrastructure/logger"
)



type CloudinaryFileUpload struct {
	CloudinaryConfig *cloudinary.Cloudinary
}

func (cdu *CloudinaryFileUpload) UploadSingleFile(data *multipart.FileHeader, publicID *string) (*string, error) {
	fileSrc, err := data.Open()
	if err != nil {
		return nil, err
	}
	res, err := cdu.CloudinaryConfig.Upload.Upload(context.Background(), fileSrc, uploader.UploadParams{
		PublicID: *publicID,
		Overwrite: utils.GetBooleanPointer(true),
		UniqueFilename: utils.GetBooleanPointer(true),
	})
	if err != nil {
		logger.Error(errors.New("could not complete UploadSingleFile on cloudinary"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil, err
	}
	logger.Info("file uploaded successfully by Cloudinary")
	return &res.URL, nil
}

func (cdu *CloudinaryFileUpload) DeleteSingleFile(publicID string) (error) {
	_, err := cdu.CloudinaryConfig.Admin.DeleteAssets(context.Background(), admin.DeleteAssetsParams{
		PublicIDs: api.CldAPIArray{publicID},
		Invalidate: utils.GetBooleanPointer(true),
		KeepOriginal: utils.GetBooleanPointer(false),
	})
	if err != nil {
		logger.Error(errors.New("could not complete DeleteSingleFile on cloudinary"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	logger.Info("file deleted successfully by Cloudinary")
	return nil
}