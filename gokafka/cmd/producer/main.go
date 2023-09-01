package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("transferiu", "teste", producer, []byte("transferencia2"), deliveryChan)
	go DeliveryReport(deliveryChan)

	fmt.Println("Thread end.")

	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "gokafka-kafka-1:9092",
		"delivery.timeout.ms": "10",   // if 0 -> unlimited
		"acks":                "all",  // 0, 1 or all
		"enable.idempotence":  "true", // if true, acks must be all
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}
	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada", ev.TopicPartition)
			}
		}
	}
}

// CMD kafka
// kafka-topics --create --bootstrap-server:9092 --topic=teste --partitions=3
// kafka-console-consumer --bootstrap-server=localhost:9092 --topic=teste
