package storage

import (
	"fmt"
	"mime/multipart"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	storage "github.com/supabase-community/storage-go"
)

type StorageSdk struct {
	sdk      StorageClient
	buckedId string
}

/**
 * Creates an sdk instance for supabase storage
 */
func NewStorageSdk() (*StorageSdk, error) {
	apiKey := utils.GetEnv("SUPABASE_API_KEY", "")
	projectRef := utils.GetEnv("SUPABASE_URL", "")
	bucketId := utils.GetEnv("SUPABASE_BUCKET_NAME", "cheatsheets")

	if apiKey == "" || projectRef == "" {
		return &StorageSdk{}, fmt.Errorf("Missing Supabase credentials. Supabase Storage SDK not configured")
	}

	storageURL := fmt.Sprintf("%s/storage/v1", projectRef)
	client := storage.NewClient(storageURL, apiKey, nil)

	return &StorageSdk{
		sdk:      &storageClient{client},
		buckedId: bucketId,
	}, nil
}

/**
 * Uploads a file to Supabase Storage
 * @param name string - name of the file
 * @param image multipart.File - file to be uploaded
 * @return string - public URL of the uploaded file
 * @return error
 */
func (s *StorageSdk) UploadFile(fileName string, image multipart.File) (string, error) {

	_, err := s.sdk.UploadFile(s.buckedId, fileName, image)
	if err != nil {
		return "", err
	}

	// Generate public URL
	resp := s.sdk.GetPublicUrl(s.buckedId, fileName)

	return resp.SignedURL, nil
}

/**
 * Updates a file in Supabase Storage
 * @param name string - name of the file
 * @param image multipart.File - file to be uploaded
 * @return string - public URL of the uploaded file
 * @return error
 */
func (s *StorageSdk) UpdateFileInPlace(fileName string, image multipart.File) (string, error) {

	_, err := s.sdk.UpdateFile(s.buckedId, fileName, image)
	if err != nil {
		return "", err
	}

	// Generate public URL
	resp := s.sdk.GetPublicUrl(s.buckedId, fileName)

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
 * Deletes a file from Supabase Storage
 * @param name string - name of the file
 * @return error
 */
func (s *StorageSdk) DeleteFile(fileName string) error {
	files := []string{fileName}

	_, err := s.sdk.RemoveFile(s.buckedId, files)
	if err != nil {
		return err
	}

	return nil
}
