package filemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// FileManager defines the interface for file management operations.
type FileManager interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

// LocalFileManager implements the FileManager interface for local file storage.
type LocalFileManager struct {
	uploadDir string
}

// NewLocalFileManager creates a new instance of LocalFileManager.
func NewLocalFileManager(uploadDir string) *LocalFileManager {
	return &LocalFileManager{
		uploadDir: uploadDir,
	}
}

// UploadFile uploads a file to the local file system.
func (fm *LocalFileManager) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory limit
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate unique filename
	filename := uuid.New().String() + filepath.Ext(filepath.Base(header.Filename))
	filePath := filepath.Join(fm.uploadDir, filename)

	// Ensure the upload directory exists
	err = os.MkdirAll(fm.uploadDir, 0755) // Create directory if it doesn't exist
	if err != nil {
		http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "Could not write file", http.StatusInternalServerError)
		return
	}

	// Return file path as response
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(filename)) //_, err = w.Write([]byte(filePath))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

//********************************TODO*******************************************************************************

// AWSFileManager implements the FileManager interface for AWS S3 storage.
// In order to interact with S3 we will need to use the AWS SDK for Go Language
type AWSFileManager struct {
	// ...Define your AWS S3 credentials and configuration
}

// NewAWSFileManager creates a new instance of AWSFileManager.
func NewAWSFileManager() *AWSFileManager { // ... AWS S3 credentials and configuration) *AWSFileManager {
	// ... Initialize AWS S3 client
	return &AWSFileManager{}
}

// UploadFile uploads a file to AWS S3.
func (afm *AWSFileManager) UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	// ... use the AWS SDK for Go to upload the file to S3

	return "OK", nil
}
