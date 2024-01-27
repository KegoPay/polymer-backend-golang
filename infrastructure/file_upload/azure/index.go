package azure

import (
	"fmt"
	"time"

	"kego.com/infrastructure/file_upload/types"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblob_sas "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)



type AzureBlobSignedURLService struct {
	AccountName string
	ContainerName string
	AccountKey string
}

func (azurlservice *AzureBlobSignedURLService) GeneratedSignedURL(file_name string, permission types.SignedURLPermission) (*string, error){
	credential, err := azblob.NewSharedKeyCredential(azurlservice.AccountName, azurlservice.AccountKey)
	if err != nil {
		return nil, err
	}
	
	sasQueryParams, err := azblob_sas.BlobSignatureValues{
		Protocol:      azblob_sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(10 * time.Minute), // url is valid for only 10 mins
		Permissions:   to.Ptr(azblob_sas.BlobPermissions{Read: permission.Read, Write: permission.Write}).String(),
		ContainerName: azurlservice.ContainerName,
	}.SignWithSharedKey(credential)
	if err != nil {
		return nil, err
	}
	
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", azurlservice.AccountName, sasQueryParams.Encode())
	return &sasURL, nil
}

func (azurlservice *AzureBlobSignedURLService) DeleteFileByURL(file_url string) (bool, error) {
	// blobServiceClient, _ := azblob.newcl(azurlservice.AccountName, azurlservice.AccountKey)
	// blobServiceClient.
	// containerClient := blobServiceClient.con(containerName)
	// blobClient := containerClient.NewBlobClient(blobName)
	return false, nil
}