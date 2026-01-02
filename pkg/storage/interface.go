package storage

import (
	"mime/multipart"

	supabase "github.com/supabase-community/storage-go"
)

/** Interface for Supabase Storage SDK **/
type StorageService interface {
	UploadFile(fileName string, image multipart.File) (string, error)
	UpdateFileInPlace(fileName string, image multipart.File) (string, error)
	ReplaceFile(prevFileName, newFileName string, image multipart.File) (string, error)
	DeleteFile(fileName string) error
}

/** Interface for Supabase Storage Client **/

type StorageClient interface {
	UploadFile(buckedId string, fileName string, image multipart.File) (supabase.FileUploadResponse, error)
	UpdateFile(buckedId string, fileName string, image multipart.File) (supabase.FileUploadResponse, error)
	RemoveFile(bucketId string, paths []string) ([]supabase.FileUploadResponse, error)
	GetPublicUrl(buckedId string, fileName string) supabase.SignedUrlResponse
}

type storageClient struct {
	client *supabase.Client
}

func (r *storageClient) UploadFile(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error) {
	return r.client.UploadFile(bucketId, path, file)
}

func (r *storageClient) UpdateFile(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error) {
	return r.client.UpdateFile(bucketId, path, file)
}

func (r *storageClient) RemoveFile(bucketId string, paths []string) ([]supabase.FileUploadResponse, error) {
	return r.client.RemoveFile(bucketId, paths)
}

func (r *storageClient) GetPublicUrl(bucketId string, filePath string) supabase.SignedUrlResponse {
	return r.client.GetPublicUrl(bucketId, filePath)
}
