package upload

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/fauzancodes/sales-demo-api/app/config"
)

func UploadFile(file interface{}, folder string, filename string) (secureUrl, publicID, cloudName string, err error) {
	cloudName = config.LoadConfig().CloudinaryCloudName
	apiKey := config.LoadConfig().CloudinaryAPIKey
	apiSecret := config.LoadConfig().CLoudinaryAPISecret
	folder = config.LoadConfig().CloudinaryFolder + "/" + folder
	request, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	response, err := request.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:   folder,
		PublicID: filename,
	})
	if err != nil {
		return
	}

	secureUrl = response.SecureURL
	publicID = response.PublicID

	return
}
