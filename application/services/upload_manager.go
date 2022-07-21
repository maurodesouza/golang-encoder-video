package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type UploadManager struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewUploadManager() *UploadManager {
	return &UploadManager{}
}

func (uploadManager *UploadManager) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	storage_path := os.Getenv("LOCAL_STORAGE_PATH")

	filename := strings.Split(objectPath, storage_path+"/")[1]
	file, err := os.Open(objectPath)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := client.Bucket(uploadManager.OutputBucket).Object(filename).NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	_, err = io.Copy(writer, file)

	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	return nil
}

func (uploadManager *UploadManager) loadVideoPaths() error {
	err := filepath.Walk(uploadManager.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			uploadManager.Paths = append(uploadManager.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
