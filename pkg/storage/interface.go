package storage

import (
	"mime/multipart"

	storage "github.com/supabase-community/storage-go"
)

type StorageClient interface {
	UploadFile(buckedId string, fileName string, image multipart.File) (storage.FileUploadResponse, error)
	UpdateFile(buckedId string, fileName string, image multipart.File) (storage.FileUploadResponse, error)
	RemoveFile(bucketId string, paths []string) ([]storage.FileUploadResponse, error)
	GetPublicUrl(buckedId string, fileName string) storage.SignedUrlResponse
}

type storageClient struct {
	client *storage.Client
}

func (r *storageClient) UploadFile(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
	return r.client.UploadFile(bucketId, path, file)
}

func (r *storageClient) UpdateFile(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
	return r.client.UpdateFile(bucketId, path, file)
}

func (r *storageClient) RemoveFile(bucketId string, paths []string) ([]storage.FileUploadResponse, error) {
	return r.client.RemoveFile(bucketId, paths)
}

func (r *storageClient) GetPublicUrl(bucketId string, filePath string) storage.SignedUrlResponse {
	return r.client.GetPublicUrl(bucketId, filePath)
}
