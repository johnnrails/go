package internal

import "github.com/confluentinc/confluent-kafka-go/kafka"

// infinite loop that must be send to a goroutine when calling this function
func Consume(topics []string, servers string, msgChan chan *kafka.Message) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          "imersao",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	consumer.SubscribeTopics(topics, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		}
	}
}
