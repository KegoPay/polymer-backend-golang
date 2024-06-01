package azure

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"usepolymer.co/infrastructure/file_upload/types"
	"usepolymer.co/infrastructure/logger"

	_azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblob_sas "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	azblob "github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureBlobSignedURLService struct {
	AccountName   string
	ContainerName string
	AccountKey    string
}

func (azurlservice *AzureBlobSignedURLService) GeneratedSignedURL(file_name string, permission types.SignedURLPermission) (*string, error) {
	if permission.Read == permission.Write {
		return nil, errors.New("permission must be either read or write")
	}
	_credential, err := _azblob.NewSharedKeyCredential(azurlservice.AccountName, azurlservice.AccountKey)
	if err != nil {
		logger.Error(errors.New("error generated _azblob shared key credential"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}

	credential, err := azblob.NewSharedKeyCredential(azurlservice.AccountName, azurlservice.AccountKey)
	if err != nil {
		logger.Error(errors.New("error generated azblob shared key credential"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", azurlservice.AccountName, azurlservice.ContainerName, file_name))
	if err != nil {
		logger.Error(errors.New("error parsing shared token url"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}
	blobURL := azblob.NewBlockBlobURL(*URL, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	sasQueryParams, err := azblob_sas.BlobSignatureValues{
		Protocol:      azblob_sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(5 * time.Minute), // url is valid for only 5 mins
		Permissions:   (&azblob_sas.BlobPermissions{Read: permission.Read, Write: permission.Write}).String(),
		ContainerName: azurlservice.ContainerName,
		BlobName:      file_name,
	}.SignWithSharedKey(_credential)
	if err != nil {
		logger.Error(errors.New("error blob signature values"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return nil, err
	}
	sasURL := fmt.Sprintf("%s?%s", blobURL.String(), sasQueryParams.Encode())
	return &sasURL, nil
}

// WARNING: For internal use only!!
func (azurlservice *AzureBlobSignedURLService) UploadBase64File(file_name string, file *string) error {
	credential, err := azblob.NewSharedKeyCredential(azurlservice.AccountName, azurlservice.AccountKey)
	if err != nil {
		logger.Error(errors.New("error generating shared key cred for UploadBase64File"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", azurlservice.AccountName, azurlservice.ContainerName, file_name))
	if err != nil {
		logger.Error(errors.New("error parsing url in UploadBase64File"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	blobURL := azblob.NewBlockBlobURL(*URL, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	options := azblob.UploadStreamToBlockBlobOptions{
		MaxBuffers: 3,
	}
	data, err := base64.StdEncoding.DecodeString(*file)
	if err != nil {
		logger.Error(errors.New("error converting base64 file to []byte for UploadBase64File"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	res, err := azblob.UploadStreamToBlockBlob(context.Background(), bytes.NewReader(data), blobURL, options)
	if err != nil {
		logger.Error(errors.New("error uploading base64 file to azure"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
	}
	logger.Info("file uploaded successfully", logger.LoggerOptions{
		Key:  "",
		Data: res.Response().StatusCode,
	})
	return nil
}
