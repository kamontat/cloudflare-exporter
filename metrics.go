package main

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/metrics"
	"github.com/kamontat/cloudflare-exporter/prom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func newMetricSet(client *cloudflare.Client, p *prom.Prometheus, config *viper.Viper) (metricSet *MetricSet) {
	metricSet = &MetricSet{client: client, prom: p, config: config}

	metricSet.JobExecTotal = prom.MustRegister(p, prom.SUBSYS_JOB, "exec_total",
		func(name string, config *viper.Viper) *prometheus.CounterVec {
			return prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: name,
				Help: "Number of scheduler job executed",
			}, []string{"uuid", "name", "status"})
		},
	)

	metricSet.JobExecTime = prom.MustRegister(p, prom.SUBSYS_JOB, "exec_time_seconds",
		func(name string, config *viper.Viper) *prometheus.HistogramVec {
			return prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: name,
				Help: "How long scheduler job executed",
				// Start with 100ms, increase 50% each bucket and stop at 10 buckets
				Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 10),
			}, []string{"uuid", "name"})
		},
	)

	metricSet.ZoneRequestTotal = prom.MustRegister(p, prom.SUBSYS_ZONE, "requests_total",
		func(name string, config *viper.Viper) *prometheus.CounterVec {
			return prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: name,
				Help: "Number of requests per zones",
			}, []string{"account", "zone"})
		},
	)

	return metricSet
}

type MetricSet struct {
	client *cloudflare.Client
	prom   *prom.Prometheus
	config *viper.Viper

	JobExecTotal     *prometheus.CounterVec
	JobExecTime      *prometheus.HistogramVec
	ZoneRequestTotal *prometheus.CounterVec
}

func (m *MetricSet) Fetch(ctx context.Context, scheduler gocron.Scheduler) error {
	logger.Debug("fetch metrics data from cloudflare")

	var job gocron.Job
	for _, j := range scheduler.Jobs() {
		if j.Name() == JOB_NAME_DEFAULT {
			job = j
		}
	}

	var fetcher = metrics.New(ctx, m.client, m.config, job)

	go fetcher.ZoneRequest(m.ZoneRequestTotal)

	fetcher.Wait()

	lastTime, err1 := job.LastRun()
	nextTime, err2 := job.NextRun()
	if err1 == nil && err2 == nil {
		logger.Debug(fmt.Sprintf("Next time will be %s (last %s)",
			nextTime.String(),
			lastTime.String(),
		))
	}

	return nil
}

func (m *MetricSet) JobMonitor() *prom.JobMonitor {
	return prom.NewJobMonitor(logger, m.JobExecTotal, m.JobExecTime)
}
