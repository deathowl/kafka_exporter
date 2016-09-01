package main

import (
        "strconv"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/Shopify/sarama"
)


type kafkaCollector struct {
	upIndicator *prometheus.Desc
	topicsCount *prometheus.Desc
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
		topicsCount: prometheus.NewDesc("kafka_topics", "Count of topics", nil, nil),
		metrics: map[string]kafkaMetric{
			"kafka_topics": {
				desc:    prometheus.NewDesc("kafka_placeholder", "Placeholder metric", nil, nil),
				extract: func(topics string) float64 { return parseFloatOrZero(topics) },
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

	var topics []string
        var kafkaConsumer sarama.Consumer
	var err error
        kafkaConsumer, err = sarama.NewConsumer(strings.Split(*kafkaAddr, ","), nil)
	if err != nil {
		log.Fatalln("Failed to start Sarama consumer:", err)
	}
	
	topics, err = kafkaConsumer.Topics()
	if err !=nil {
		log.Error("Failed to get list of topics")
		ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 0)
	}
	ch <- prometheus.MustNewConstMetric(c.topicsCount, prometheus.GaugeValue, float64(len(topics)))
	log.Info(topics)
	ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 1)
}
