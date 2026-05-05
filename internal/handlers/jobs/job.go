package jobHandler

// import (
// 	"time"

// 	"go.uber.org/zap"
// )

// type Job struct {
// 	Name        string
// 	Frequency   int // in seconds
// 	HandlerFunc JobFunc
// 	RunOnInit   bool
// }

// type JobFunc func(name string, registery services.ServiceRegistry, logger *zap.Logger)

// func NewJob(name string, frequency int, handlerFunc JobFunc, runOnInit bool) Job {
// 	return Job{
// 		Name:        name,
// 		Frequency:   frequency,
// 		HandlerFunc: handlerFunc,
// 		RunOnInit:   runOnInit,
// 	}
// }

// type JobsHandler interface {
// 	RegisterJobs(job Job)
// 	RunJobs()
// }

// func NewJobsHandler(serviceRegistry services.ServiceRegistry, logger *zap.Logger) JobsHandler {
// 	return &jobsHandler{
// 		serviceRegistry: serviceRegistry,
// 		logger:          logger,
// 		jobs:            make([]Job, 0),
// 	}
// }

// type jobsHandler struct {
// 	serviceRegistry services.ServiceRegistry
// 	logger          *zap.Logger
// 	jobs            []Job
// }

// func (j *jobsHandler) RegisterJobs(job Job) {
// 	j.jobs = append(j.jobs, job)
// }

// func (j *jobsHandler) RunJobs() {
// 	for _, job := range j.jobs {
// 		go j.executeJob(job)
// 		j.logger.Info("Running  job: " + job.Name)
// 	}
// }

// func (j *jobsHandler) executeJob(job Job) {
// 	j.logger.Info("Executing job: " + job.Name)
// 	defer func() {
// 		j.logger.Info("Error executing job: " + job.Name)
// 		if r := recover(); r != nil {
// 			j.logger.Error("Recovered from panic in job: "+job.Name, zap.Any("error", r))
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 		go j.executeJob(job)
// 	}()
// 	ticker := time.NewTicker(time.Duration(job.Frequency) * time.Second)
// 	defer ticker.Stop()
// 	done := make(chan bool)

// 	if job.RunOnInit {
// 		job.HandlerFunc(job.Name, j.serviceRegistry, j.logger)
// 	}

// 	for {
// 		select {
// 		case <-done:
// 			return
// 		case <-ticker.C:
// 			job.HandlerFunc(job.Name, j.serviceRegistry, j.logger)
// 		}
// 	}
// }
