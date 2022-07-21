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

func (service *VideoService) Encode() error {
	storage_path, video_id := os.Getenv("LOCAL_STORAGE_PATH"), service.Video.Id
	base_path := storage_path + "/" + video_id

	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, base_path+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, base_path)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	command := exec.Command("mp4dash", cmdArgs...)
	output, err := command.CombinedOutput()

	if err != nil {
		return err
	}

	if len(output) > 0 {
		log.Printf("============> output: %s\n", string(output))
	}

	return nil
}

func (service *VideoService) Finish() error {
	storage_path, video_id := os.Getenv("LOCAL_STORAGE_PATH"), service.Video.Id
	base_path := storage_path + "/" + video_id

	err := os.Remove(base_path + ".mp4")

	if err != nil {
		log.Println("error removing mp4 ", video_id+".mp4")
		return err
	}

	err = os.Remove(base_path + ".frag")

	if err != nil {
		log.Println("error removing mp4 ", video_id+".frag")
		return err
	}

	err = os.RemoveAll(base_path)

	if err != nil {
		log.Println("error removing mp4 ", video_id+".mp4")
		return err
	}

	log.Println("all files have been removed: ", video_id)

	return nil
}
