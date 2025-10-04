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
		return &StorageSdk{}, fmt.Errorf("Missing Supabase credentials. Supabase Storage SDK not configured")
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

/**
 * Update a file in Supabase Storage
 * @param name string - name of the file
 * @param image multipart.File - file to be uploaded
 * @return string - public URL of the uploaded file
 * @return error
 */
func (s *StorageSdk) UpdateFileInPlace(fileName string, image multipart.File) (string, error) {

	_, err := s.client.UpdateFile(s.buckedId, fileName, image)
	if err != nil {
		return "", err
	}

	// Generate public URL
	resp := s.client.GetPublicUrl(s.buckedId, fileName)

	return resp.SignedURL, nil

}

/**
 * Removes previous image and uploads a new one
 * @param name string - previous file name
 * @param name string - new file name
 * @return string - public URL of the uploaded file
 */
func (s *StorageSdk) ReplaceFile(prevFileName, newFileName string, image multipart.File) (string, error) {
	// remove previous file
	if err := s.DeleteFile(prevFileName); err != nil {
		return "", err
	}

	// upload new image in new location
	url, err := s.UploadFile(newFileName, image)
	if err != nil {
		return "", err
	}

	return url, nil

}

/**
 * Delete a file from Supabase Storage
 * @param name string - name of the file
 * @return error
 */
func (s *StorageSdk) DeleteFile(fileName string) error {
	files := []string{fileName}

	_, err := s.client.RemoveFile(s.buckedId, files)
	if err != nil {
		return err
	}

	return nil
}
