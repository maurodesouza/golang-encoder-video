package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepositoryInterface
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (service *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(service.Video.FilePath)

	reader, err := object.NewReader(ctx)

	if err != nil {
		return err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	storage_path := os.Getenv("LOCAL_STORAGE_PATH")
	file, err := os.Create(storage_path + "/" + service.Video.Id + ".mp4")

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	defer file.Close()

	log.Printf("video %v has been stored", service.Video.Id)

	return nil
}

func (service *VideoService) Fragment() error {
	storage_path, video_id := os.Getenv("LOCAL_STORAGE_PATH"), service.Video.Id
	base_path := storage_path + "/" + video_id

	err := os.Mkdir(base_path, os.ModePerm)

	if err != nil {
		return err
	}

	source := base_path + ".mp4"
	target := base_path + ".frag"

	command := exec.Command("mp4fragment", source, target)
	output, err := command.CombinedOutput()

	if err != nil {
		return err
	}

	if len(output) > 0 {
		log.Printf("============> output: %s\n", string(output))
	}

	return nil
}
