package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepositoryInterface
	VideoService  VideoService
}

func (jobService *JobService) Start() error {
	err := jobService.changeJobStatus("DOWNLOADING")

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.VideoService.Download(os.Getenv("INPUT_BUCKET_NAME"))

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.changeJobStatus("FRAGMENTING")

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.VideoService.Fragment()

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.changeJobStatus("ENCODING")

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.VideoService.Encode()

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.performUpload()

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.changeJobStatus("FINISHING")

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.VideoService.Finish()

	if err != nil {
		return jobService.failJob(err)
	}

	err = jobService.changeJobStatus("COMPLETED")

	if err != nil {
		return jobService.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {

	err := j.changeJobStatus("UPLOADING")

	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewUploadManager()
	videoUpload.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + j.VideoService.Video.Id

	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err
}

func (jobService *JobService) changeJobStatus(status string) error {
	var err error

	jobService.Job.Status = status
	jobService.Job, err = jobService.JobRepository.Update(jobService.Job)

	if err != nil {
		return jobService.failJob(err)
	}

	return nil
}

func (jobService *JobService) failJob(error error) error {
	jobService.Job.Status = "FAILED"
	jobService.Job.Error = error.Error()

	_, err := jobService.JobRepository.Update(jobService.Job)

	if err != nil {
		return err
	}

	return error
}
