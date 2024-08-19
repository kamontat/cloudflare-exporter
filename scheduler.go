package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/kamontat/cloudflare-exporter/utils"
	"go.uber.org/zap"
)

var (
	jobIds map[string]uuid.UUID

	JOB_NAME_DEFAULT = "default"
)

func init() {
	jobIds = make(map[string]uuid.UUID)
	// UUID V4
	jobIds[JOB_NAME_DEFAULT] = uuid.MustParse("c1465b30-ea0e-4a62-a5ef-1db08fedc9e1")
}

func newScheduler(metrics *MetricSet) gocron.Scheduler {
	// Add 5 seconds offset from cloudflare timeout
	var timeoutOffset = 5 * time.Second
	return utils.CheckErrorWithData(gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithStopTimeout(config.GetDuration(configs.CONF_CF_TIMEOUT)+timeoutOffset),
		gocron.WithMonitor(metrics.JobMonitor()),
		gocron.WithLogger(loggers.GocronLoggerAdapter(logger)),
		gocron.WithLimitConcurrentJobs(1, gocron.LimitModeReschedule),
		gocron.WithGlobalJobOptions(
			gocron.WithStartAt(gocron.WithStartImmediately()),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					logger.Debug(fmt.Sprintf("Start run job %s (%s)", jobName, jobID))
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					logger.Debug(fmt.Sprintf("Stop run job %s (%s)", jobName, jobID))
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					logger.Error(fmt.Sprintf("Job %s (%s) stopped with error", jobName, jobID), zap.Error(err))
				}),
				gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
					logger.Error(fmt.Sprintf("Job %s (%s) stopped with panic", jobName, jobID), zap.Any("recover", recoverData))
				}),
			),
		),
	))
}

func setupJob(scheduler gocron.Scheduler, name string, def gocron.JobDefinition, executor func(context.Context, gocron.Scheduler) error) {
	jid := jobIds[name]
	job := utils.CheckErrorWithData(scheduler.NewJob(
		def, gocron.NewTask(executor, context.Background(), scheduler),
		gocron.WithIdentifier(jid),
		gocron.WithName(name),
	))

	logger.Debug(fmt.Sprintf("Set up Job %s", job.Name()), zap.String("id", job.ID().String()))
}

func startScheduler(scheduler gocron.Scheduler) {
	scheduler.Start()
}

func shutdownScheduler(scheduler gocron.Scheduler) {
	logger.Info("Shutdown cronjob scheduler", zap.Error(scheduler.Shutdown()))
}
