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
	partitionsCount * prometheus.Desc
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
		partitionsCount: prometheus.NewDesc("kafka_partitions", "Count of partitions per topic",  []string{"topic"}, nil),
	}
}

func (c *kafkaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upIndicator
	ch <- c.topicsCount
	ch <- c.partitionsCount
}

func (c *kafkaCollector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Fetching metrics from Kafka")

	var topics []string
        var kafkaConsumer sarama.Consumer
	var err error
	var partitions []int32
        kafkaConsumer, err = sarama.NewConsumer(strings.Split(*kafkaAddr, ","), nil)
	if err != nil {
		log.Fatalln("Failed to start Sarama consumer:", err)
	}
	
	ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 1)
	topics, err = kafkaConsumer.Topics()
	if err !=nil {
		log.Error("Failed to get list of topics")
		ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 0)
	}
	ch <- prometheus.MustNewConstMetric(c.topicsCount, prometheus.GaugeValue, float64(len(topics)))
	for _, topic := range topics{
		partitions,_ = kafkaConsumer.Partitions(topic)
		ch <- prometheus.MustNewConstMetric(c.partitionsCount, prometheus.GaugeValue,float64(len(partitions)), topic)

	}
}
