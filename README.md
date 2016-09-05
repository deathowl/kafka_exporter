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

**Before first run**

Dont forget to add consumer group, and allow acces to it via kafka ACL
 ./kafka-acls.sh --authorizer kafka.security.auth.SimpleAclAuthorizer --authorizer-properties zookeeper.connect=zookeeper.host:2181 --allow-principal User:* --operation All --topic test --group zkexporter.om --add 
