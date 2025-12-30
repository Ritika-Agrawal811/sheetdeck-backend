package storage

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"testing"

	storage "github.com/supabase-community/storage-go"
)

/** Mock Implementation for Supabase Storage **/

/**
 * mockStorageClient implements StorageClient interface
 * Each field is a function from the interface
 * This helps us to customize behavior per test case
 */
type mockStorageClient struct {
	uploadFunc    func(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error)
	updateFunc    func(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error)
	removeFunc    func(bucketId string, paths []string) ([]storage.FileUploadResponse, error)
	getPublicFunc func(bucketId string, filePath string) storage.SignedUrlResponse
}

/* Mocks UploadFile */
func (m *mockStorageClient) UploadFile(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(bucketId, path, file)
	}

	return storage.FileUploadResponse{Key: path}, nil
}

/* Mocks UpdateFile */
func (m *mockStorageClient) UpdateFile(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(bucketId, path, file)
	}
	return storage.FileUploadResponse{Key: path}, nil
}

/* Mocks RemoveFile */
func (m *mockStorageClient) RemoveFile(bucketId string, paths []string) ([]storage.FileUploadResponse, error) {
	if m.removeFunc != nil {
		return m.removeFunc(bucketId, paths)
	}
	return []storage.FileUploadResponse{{Key: paths[0]}}, nil
}

/* Mocks GetPublicUrl */
func (m *mockStorageClient) GetPublicUrl(bucketId string, filePath string) storage.SignedUrlResponse {
	if m.getPublicFunc != nil {
		return m.getPublicFunc(bucketId, filePath)
	}
	return storage.SignedUrlResponse{
		SignedURL: fmt.Sprintf("https://example.supabase.co/storage/v1/object/public/%s/%s", bucketId, filePath),
	}
}

/** Mock Implementation for multipart.File **/

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func createMockFile(content string) multipart.File {
	return &mockMultipartFile{
		Reader: bytes.NewReader([]byte(content)),
	}
}

/** Tests for Supabase Storage SDK **/
func TestNewStorageSdk(t *testing.T) {
	tests := []struct {
		name        string
		projectRef  string
		apiKey      string
		buckedId    string
		setEnv      bool
		expectError bool
	}{
		{
			name:        "creates SDK with valid credentials",
			projectRef:  "test-project-ref",
			apiKey:      "test-api-key",
			buckedId:    "test-bucket",
			setEnv:      true,
			expectError: false,
		},
		{
			name:        "returns error when project ref is missing",
			projectRef:  "",
			apiKey:      "test-api-key",
			buckedId:    "test-bucket",
			setEnv:      true,
			expectError: true,
		},
		{
			name:        "returns error when API key is missing",
			projectRef:  "test-project-ref",
			apiKey:      "",
			buckedId:    "test-bucket",
			setEnv:      true,
			expectError: true,
		},
		{
			name:        "uses default bucket name when not provided",
			projectRef:  "test-project-ref",
			apiKey:      "test-api-key",
			buckedId:    "",
			setEnv:      true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup env variables
			if tt.setEnv {
				if tt.apiKey != "" {
					t.Setenv("SUPABASE_API_KEY", tt.apiKey)
				}

				if tt.projectRef != "" {
					t.Setenv("SUPABASE_URL", tt.apiKey)
				}

				if tt.buckedId != "" {
					t.Setenv("SUPABASE_BUCKET_NAME", tt.buckedId)
				}
			}

			storageClient, err := NewStorageSdk()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if storageClient.sdk == nil {
					t.Error("Expected client to be initialized")
				}

				/* Check bucket name */
				expectedBucket := tt.buckedId
				if expectedBucket == "" {
					expectedBucket = "cheatsheets" // default value
				}
				if storageClient.buckedId != expectedBucket {
					t.Errorf("Expected bucketId %q, got %q", expectedBucket, storageClient.buckedId)
				}
			}
		})
	}
}

func TestUploadFile(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		bucketId    string
		uploadError error
		expectedURL string
		expectError bool
	}{
		{
			name:        "successfully uploads file",
			fileName:    "test.png",
			bucketId:    "cheatsheets",
			uploadError: nil,
			expectedURL: "https://example.supabase.co/storage/v1/object/public/cheatsheets/test.png",
			expectError: false,
		},
		{
			name:        "handles upload error",
			fileName:    "test.png",
			bucketId:    "cheatsheets",
			uploadError: fmt.Errorf("network error"),
			expectError: true,
		},
		{
			name:        "uploads file with path",
			fileName:    "images/test.png",
			bucketId:    "cheatsheets",
			uploadError: nil,
			expectedURL: "https://example.supabase.co/storage/v1/object/public/cheatsheets/images/test.png",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates mock client */
			mockClient := &mockStorageClient{
				uploadFunc: func(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
					/* Verify correct bucket is used */
					if bucketId != tt.bucketId {
						t.Errorf("Expected bucketId %q, but got %q", tt.bucketId, bucketId)
					}

					/* Verify correct path */
					if path != tt.fileName {
						t.Errorf("Expected path %q, got %q", tt.fileName, path)
					}

					if tt.uploadError != nil {
						return storage.FileUploadResponse{}, tt.uploadError
					}

					return storage.FileUploadResponse{Key: path}, nil
				},
			}

			/* Create SDK with mock client */
			storageClient := &StorageSdk{
				sdk:      mockClient,
				buckedId: tt.bucketId,
			}

			mockFile := createMockFile("test content")
			url, err := storageClient.UploadFile(tt.fileName, mockFile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if url != tt.expectedURL {
					t.Errorf("Expected URL %q, got %q", tt.expectedURL, url)
				}
			}
		})
	}
}

func TestUpdateFileInPlace(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		bucketId    string
		updateError error
		expectedURL string
		expectError bool
	}{
		{
			name:        "successfully updates file",
			fileName:    "test.png",
			bucketId:    "cheatsheets",
			updateError: nil,
			expectedURL: "https://example.supabase.co/storage/v1/object/public/cheatsheets/test.png",
			expectError: false,
		},
		{
			name:        "handles update error",
			fileName:    "nonexistent.png",
			bucketId:    "cheatsheets",
			updateError: fmt.Errorf("file not found"),
			expectError: true,
		},
		{
			name:        "updates file with path",
			fileName:    "images/test.png",
			bucketId:    "cheatsheets",
			updateError: nil,
			expectedURL: "https://example.supabase.co/storage/v1/object/public/cheatsheets/images/test.png",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates mock client */
			mockClient := &mockStorageClient{
				updateFunc: func(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
					/* Verify correct bucket is used */
					if bucketId != tt.bucketId {
						t.Errorf("Expected bucketId %q, but got %q", tt.bucketId, bucketId)
					}

					/* Verify correct path */
					if path != tt.fileName {
						t.Errorf("Expected path %q, but got %q", tt.fileName, path)
					}

					if tt.updateError != nil {
						return storage.FileUploadResponse{}, tt.updateError
					}

					return storage.FileUploadResponse{Key: path}, nil
				},
			}

			/* Create SDK with mock client */
			storageClient := &StorageSdk{
				sdk:      mockClient,
				buckedId: tt.bucketId,
			}

			mockFile := createMockFile("test content")
			url, err := storageClient.UpdateFileInPlace(tt.fileName, mockFile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if url != tt.expectedURL {
					t.Errorf("Expected URL %q, but got %q", tt.expectedURL, url)
				}
			}
		})
	}
}

func TestReplaceFile(t *testing.T) {
	tests := []struct {
		name         string
		prevFileName string
		newFileName  string
		bucketId     string
		deleteError  error
		uploadError  error
		expectedURL  string
		expectError  bool
	}{
		{
			name:         "successfully replaces file",
			prevFileName: "old.png",
			newFileName:  "new.png",
			bucketId:     "cheatsheets",
			deleteError:  nil,
			uploadError:  nil,
			expectedURL:  "https://example.supabase.co/storage/v1/object/public/cheatsheets/new.png",
			expectError:  false,
		},
		{
			name:         "handles delete error",
			prevFileName: "old.png",
			newFileName:  "new.png",
			bucketId:     "cheatsheets",
			deleteError:  fmt.Errorf("file not found"),
			expectError:  true,
		},
		{
			name:         "handles upload error after successful delete",
			prevFileName: "old.png",
			newFileName:  "new.png",
			bucketId:     "cheatsheets",
			deleteError:  nil,
			uploadError:  fmt.Errorf("storage full"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteWasCalled := false
			uploadWasCalled := false

			mockClient := &mockStorageClient{
				removeFunc: func(bucketId string, paths []string) ([]storage.FileUploadResponse, error) {
					deleteWasCalled = true

					/* Verify correct bucket is used */
					if bucketId != tt.bucketId {
						t.Errorf("Expected bucketId %q, but got %q", tt.bucketId, bucketId)
					}

					/* Verify correct path */
					if paths[0] != tt.prevFileName {
						t.Errorf("Expected to delete %q, but got %q", tt.prevFileName, paths[0])
					}

					if tt.deleteError != nil {
						return nil, tt.deleteError
					}

					return []storage.FileUploadResponse{{Key: paths[0]}}, nil
				},

				uploadFunc: func(bucketId string, path string, file multipart.File) (storage.FileUploadResponse, error) {
					uploadWasCalled = true

					/* Verify correct bucket is used */
					if bucketId != tt.bucketId {
						t.Errorf("Expected bucketId %q, but got %q", tt.bucketId, bucketId)
					}

					/* Verify correct path */
					if path != tt.newFileName {
						t.Errorf("Expected path %q, but got %q", tt.newFileName, path)
					}

					if tt.uploadError != nil {
						return storage.FileUploadResponse{}, tt.uploadError
					}

					return storage.FileUploadResponse{Key: path}, nil
				},
			}

			/* Create SDK with mock client */
			storageClient := &StorageSdk{
				sdk:      mockClient,
				buckedId: tt.bucketId,
			}

			mockFile := createMockFile("new content")
			url, err := storageClient.ReplaceFile(tt.prevFileName, tt.newFileName, mockFile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if url != tt.expectedURL {
					t.Errorf("Expected URL %q, got %q", tt.expectedURL, url)
				}
			}

			/* Verify correct sequence of operations */
			if !tt.expectError {
				if !deleteWasCalled {
					t.Error("Expected delete to be called")
				}

				if !uploadWasCalled {
					t.Error("Expected upload to be called")
				}
			}

			/* If delete fails, upload should not be called */
			if tt.deleteError != nil && uploadWasCalled {
				t.Error("Upload should not be called when delete fails")
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		bucketId    string
		removeError error
		expectError bool
	}{
		{
			name:        "successfully deletes file",
			fileName:    "test.png",
			bucketId:    "cheatsheets",
			removeError: nil,
			expectError: false,
		},
		{
			name:        "handles delete error of file not found",
			fileName:    "nonexistent.png",
			bucketId:    "cheatsheets",
			removeError: fmt.Errorf("file not found"),
			expectError: true,
		},
		{
			name:        "handles delete error of permission denied",
			fileName:    "protected.png",
			bucketId:    "cheatsheets",
			removeError: fmt.Errorf("permission denied"),
			expectError: true,
		},
		{
			name:        "deletes file with path",
			fileName:    "images/old/test.png",
			bucketId:    "cheatsheets",
			removeError: nil,
			expectError: false,
		},
		{
			name:        "handles network error",
			fileName:    "test.png",
			bucketId:    "cheatsheets",
			removeError: fmt.Errorf("network timeout"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* Creates mock client */
			mockClient := &mockStorageClient{
				removeFunc: func(bucketId string, paths []string) ([]storage.FileUploadResponse, error) {
					/* Verify correct bucket is used */
					if bucketId != tt.bucketId {
						t.Errorf("Expected bucketId %q, but got %q", tt.bucketId, bucketId)
					}

					/* Verify it deletes exactly one file */
					if len(paths) != 1 {
						t.Errorf("Expected 1 path, got %d", len(paths))
					}

					/* Verify correct file name */
					if paths[0] != tt.fileName {
						t.Errorf("Expected path %q, got %q", tt.fileName, paths[0])
					}

					if tt.removeError != nil {
						return nil, tt.removeError
					}

					return []storage.FileUploadResponse{{Key: paths[0]}}, nil
				},
			}

			/* Create SDK with mock client */
			storageClient := &StorageSdk{
				sdk:      mockClient,
				buckedId: tt.bucketId,
			}

			err := storageClient.DeleteFile(tt.fileName)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
