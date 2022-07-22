package services

import (
	"encoder/application/repositories"
	"encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepositoryInterface
	VideoService  VideoService
}

func (jobService *JobService) Start() error {
	return nil
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
