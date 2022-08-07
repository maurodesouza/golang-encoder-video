package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func (uploadManager *UploadManager) ProcessUpload(concurrency int, doneUpload chan string) error {
	input := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := uploadManager.loadVideoPaths()

	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()

	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go uploadManager.uploadWorker(input, returnChannel, uploadClient, ctx)
	}

	go func() {
		for index := 0; index < len(uploadManager.Paths); index++ {
			input <- index
		}

		close(input)
	}()

	for response := range returnChannel {
		if response != "" {
			doneUpload <- response
			break
		}
	}

	return nil
}

func (uploadManager *UploadManager) uploadWorker(input chan int, returnChannel chan string, uploadClient *storage.Client, ctx context.Context) {
	for index := range input {
		path := uploadManager.Paths[index]

		err := uploadManager.UploadObject(path, uploadClient, ctx)

		if err != nil {
			uploadManager.Errors = append(uploadManager.Errors, path)
			log.Printf("error during the upload: %v. Error: %v", path, err)
			returnChannel <- err.Error()
		}

		returnChannel <- ""
	}

	returnChannel <- "upload completed"

}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
