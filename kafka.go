package main

import (
        "strconv"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

type kafkaCollector struct {
	upIndicator *prometheus.Desc
	metrics     map[string]kafkaMetric
}
type kafkaMetric struct {
	desc          *prometheus.Desc
	extract       func(string) float64
	extractLabels func(s string) []string
	valType       prometheus.ValueType
}

func init() {
	prometheus.MustRegister(NewKafkaCollector())
}

func parseFloatOrZero(s string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Warningf("Failed to parse to float64: %s", err)
		return 0.0
	}
	return res
}
func NewKafkaCollector() *kafkaCollector {
	return &kafkaCollector{
		upIndicator: prometheus.NewDesc("kafka_up", "Exporter successful", nil, nil),
		metrics: map[string]kafkaMetric{
			"kafka_metric_1": {
				desc:    prometheus.NewDesc("zk_avg_latency", "Average latency of requests", nil, nil),
				extract: func(s string) float64 { return parseFloatOrZero(s) },
				valType: prometheus.GaugeValue,
			},
		},
	}
}

func (c *kafkaCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Debugf("Sending %d metrics descriptions", len(c.metrics))
	for _, i := range c.metrics {
		ch <- i.desc
	}
}

func (c *kafkaCollector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Fetching metrics from Kafka")


	ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 1)

}
