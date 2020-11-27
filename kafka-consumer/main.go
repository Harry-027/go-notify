package main

import (
	"context"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/kafka-consumer/mailer"
	"github.com/segmentio/kafka-go"
	"log"
)

const (
	BROKER   = "BROKER"
	GROUPID  = "GROUP_ID"
	TOPIC    = "TOPIC"
	MINBYTES = 10e3
	MAXBYTES = 10e6
)

func main() {

	// load env variables
	config.LoadConfig()
	sender := mailer.NewSender()
	broker := config.GetConfig(BROKER)
	grpId := config.GetConfig(GROUPID)
	topic := config.GetConfig(TOPIC)

	// Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  grpId,
		Topic:    topic,
		MinBytes: MINBYTES, // 10KB
		MaxBytes: MAXBYTES, // 10MB
	})

	log.Println("Consumer started listening ... ")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		receiptId, _ := sender.Send(m.Value)
		log.Println("Msg sent receiptId: ", receiptId)
	}

	if err := r.Close(); err != nil {
		log.Fatal("Failed to close reader:", err)
	}
}
