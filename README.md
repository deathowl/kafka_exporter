** Prometheus Kafka Exporter**

*Exposes basic kafka metrics to prometheus*

Args:

* --bind-addr default:9169  bind address for the metrics server
* --metrics-path  default:/metrics, path to metrics endpoint
* --kafka  default: localhost:9092 "host:port, host2:port2"  list of kafka brokers
* --log-level default :info log level

*Exposed Metrics*

* kafka_messages : Number of messages/Topic/Partition

* kafka_partitions : Number of partitions per Topic
* kafka_topics:  Number of Topics
* kafka_up  Status of Kafka Server


