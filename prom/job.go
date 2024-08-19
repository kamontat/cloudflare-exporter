package prom

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/kamontat/cloudflare-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func NewJobMonitor(logger *zap.Logger, execTotal *prometheus.CounterVec, execTime *prometheus.HistogramVec) *JobMonitor {
	return &JobMonitor{logger: logger, execTotal: execTotal, execTime: execTime}
}

type JobMonitor struct {
	logger    *zap.Logger
	execTotal *prometheus.CounterVec
	execTime  *prometheus.HistogramVec
}

func (m *JobMonitor) IncrementJob(id uuid.UUID, name string, tags []string, status gocron.JobStatus) {
	m.logger.Debug("Increment job",
		zap.String("uuid", id.String()),
		zap.String("name", name),
		zap.Strings("tags", tags),
		zap.String("status", string(status)),
	)

	utils.SafeCall(m.execTotal, func(metrics *prometheus.CounterVec) {
		metrics.WithLabelValues(id.String(), name, string(status)).Inc()
	})
}

func (m *JobMonitor) RecordJobTiming(startTime, endTime time.Time, id uuid.UUID, name string, tags []string) {
	var diff = endTime.Sub(startTime)

	m.logger.Debug("Record job timing",
		zap.String("uuid", id.String()),
		zap.String("name", name),
		zap.Strings("tags", tags),
		zap.String("execute-time", diff.String()),
	)

	utils.SafeCall(m.execTime, func(metrics *prometheus.HistogramVec) {
		metrics.WithLabelValues(id.String(), name).Observe(diff.Seconds())
	})
}
