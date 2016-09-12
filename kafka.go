package main

import (
        "strconv"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/Shopify/sarama"
        "github.com/deathowl/go-metrics-prometheus"
	"time"

)
var saramaConfig *sarama.Config

type kafkaCollector struct {
	upIndicator *prometheus.Desc
	topicsCount *prometheus.Desc
	partitionsCount * prometheus.Desc
	messageCount *prometheus.Desc
}
func init() {
	prometheus.MustRegister(NewKafkaCollector())
	saramaConfig = sarama.NewConfig()
	pClient := prometheusmetrics.NewPrometheusProvider(saramaConfig.MetricRegistry, "kafka", "broker", prometheus.DefaultRegisterer, 1*time.Second)
	go pClient.UpdatePrometheusMetrics()
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
		messageCount: prometheus.NewDesc("kafka_messages", "Count of messages per topic per partition",  []string{"topic", "partition"}, nil),
	}
}

func (c *kafkaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upIndicator
	ch <- c.topicsCount
	ch <- c.partitionsCount
	ch <- c.messageCount
}

func (c *kafkaCollector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Fetching metrics from Kafka")

	kClient, _ := sarama.NewClient(strings.Split(*kafkaAddr, ","), saramaConfig)
	kafkaConsumer, err := sarama.NewConsumer(strings.Split(*kafkaAddr, ","), nil)
	if err != nil {
		log.Error("Failed to start Sarama consumer:", err)
		ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 1)
	topics, err := kafkaConsumer.Topics()
	if err !=nil {
		log.Error("Failed to get list of topics")
		ch <- prometheus.MustNewConstMetric(c.upIndicator, prometheus.GaugeValue, 0)
	}
	ch <- prometheus.MustNewConstMetric(c.topicsCount, prometheus.GaugeValue, float64(len(topics)))
	for _, topic := range topics{
		if topic == "__consumer_offsets"{
			continue
		}
		partitions, _ := kafkaConsumer.Partitions(topic)
		ch <- prometheus.MustNewConstMetric(c.partitionsCount, prometheus.GaugeValue,float64(len(partitions)), topic)

		for _, partition := range partitions {
			offsetManager, _ :=sarama.NewOffsetManagerFromClient("zkexporter.om", kClient)
			pom, _ := offsetManager.ManagePartition(topic, partition)
			oldOffset, _ := pom.NextOffset()
			newOffset, _ := kClient.GetOffset(topic, partition, sarama.OffsetNewest)
			ch <- prometheus.MustNewConstMetric(c.messageCount, prometheus.GaugeValue,float64(newOffset - oldOffset), topic, strconv.Itoa(int(partition)))
			pom.MarkOffset(newOffset, "already reported")
			pom.Close()
			offsetManager.Close()
		}
	}
	kafkaConsumer.Close()
	kClient.Close()
}
