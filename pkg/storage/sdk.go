package storage

import (
	"fmt"
	"mime/multipart"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	storage "github.com/supabase-community/storage-go"
)

type StorageSdk struct {
	client   *storage.Client
	buckedId string
}

func NewStorageSdk() (*StorageSdk, error) {
	apiKey := utils.GetEnv("SUPABASE_API_KEY", "")
	projectRef := utils.GetEnv("SUPABASE_URL", "")
	bucketId := utils.GetEnv("SUPABASE_BUCKET_NAME", "cheatsheets")

	// Check if the SDK is configured
	if apiKey == "" || projectRef == "" {
		return &StorageSdk{}, fmt.Errorf("Miissing Supabase credentials. Supabase Storage SDK not configured")
	}

	// Initialize the Supabase Storage client
	storageURL := fmt.Sprintf("%s/storage/v1", projectRef)
	client := storage.NewClient(storageURL, apiKey, nil)

	return &StorageSdk{
		client:   client,
		buckedId: bucketId,
	}, nil
}

/**
 * Upload a file to Supabase Storage
 * @param name string - name of the file
 * @param image multipart.File - file to be uploaded
 * @return string - public URL of the uploaded file
 * @return error
 */
func (s *StorageSdk) UploadFile(fileName string, image multipart.File) (string, error) {

	_, err := s.client.UploadFile(s.buckedId, fileName, image)
	if err != nil {
		return "", err
	}

	// Generate public URL
	resp := s.client.GetPublicUrl(s.buckedId, fileName)

	return resp.SignedURL, nil
}
