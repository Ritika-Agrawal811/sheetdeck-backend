package mocks

import (
	"bytes"
	"fmt"
	"mime/multipart"

	supabase "github.com/supabase-community/storage-go"
)

/** Mock Implementation for Supabase Storage SDK **/

type MockStorageService struct {
	UploadFileFunc        func(fileName string, image multipart.File) (string, error)
	UpdateFileInPlaceFunc func(fileName string, image multipart.File) (string, error)
	ReplaceFileFunc       func(prevFileName, newFileName string, image multipart.File) (string, error)
	DeleteFileFunc        func(fileName string) error
}

func (m *MockStorageService) UploadFile(fileName string, image multipart.File) (string, error) {
	if m.UploadFileFunc != nil {

		return m.UploadFileFunc(fileName, image)
	}

	return "https://example.com/" + fileName, nil
}

func (m *MockStorageService) UpdateFileInPlace(fileName string, image multipart.File) (string, error) {
	if m.UpdateFileInPlaceFunc != nil {

		return m.UpdateFileInPlaceFunc(fileName, image)
	}

	return "https://example.com/" + fileName, nil
}

func (m *MockStorageService) ReplaceFile(prevFileName, newFileName string, image multipart.File) (string, error) {
	if m.ReplaceFileFunc != nil {

		return m.ReplaceFileFunc(prevFileName, newFileName, image)
	}

	return "https://example.com/" + newFileName, nil
}

func (m *MockStorageService) DeleteFile(fileName string) error {
	if m.DeleteFileFunc != nil {

		return m.DeleteFileFunc(fileName)
	}

	return nil
}

/** Mock Implementation for Supabase Storage Client **/

type MockStorageClient struct {
	UploadFunc       func(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error)
	UpdateFunc       func(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error)
	RemoveFunc       func(bucketId string, paths []string) ([]supabase.FileUploadResponse, error)
	GetPublicUrlFunc func(bucketId string, filePath string) supabase.SignedUrlResponse
}

func (m *MockStorageClient) UploadFile(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(bucketId, path, file)
	}

	return supabase.FileUploadResponse{Key: path}, nil
}

func (m *MockStorageClient) UpdateFile(bucketId string, path string, file multipart.File) (supabase.FileUploadResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(bucketId, path, file)
	}
	return supabase.FileUploadResponse{Key: path}, nil
}

func (m *MockStorageClient) RemoveFile(bucketId string, paths []string) ([]supabase.FileUploadResponse, error) {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(bucketId, paths)
	}
	return []supabase.FileUploadResponse{{Key: paths[0]}}, nil
}

func (m *MockStorageClient) GetPublicUrl(bucketId string, filePath string) supabase.SignedUrlResponse {
	if m.GetPublicUrlFunc != nil {
		return m.GetPublicUrlFunc(bucketId, filePath)
	}
	return supabase.SignedUrlResponse{
		SignedURL: fmt.Sprintf("https://example.supabase.co/storage/v1/object/public/%s/%s", bucketId, filePath),
	}
}

/** Mock Implementation for multipart.File **/
type MockMultipartFile struct {
	*bytes.Reader
}

func (m *MockMultipartFile) Close() error {
	return nil
}

func CreateMockFile(content string) multipart.File {
	return &MockMultipartFile{
		Reader: bytes.NewReader([]byte(content)),
	}
}
