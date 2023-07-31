package session

import (
	"Social-Net-Dialogs/internal/helper"
	"Social-Net-Dialogs/internal/store"
	"Social-Net-Dialogs/models"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"os"
)

type Consumer struct {
	Reader       *kafka.Reader
	SessionStore store.SessionStore
}

func NewSessionConsumer(sessionStore store.SessionStore) Consumer {
	brokerHost := helper.GetEnvValue("KAFKA_BROKER_HOST", "localhost")
	brokerPort := helper.GetEnvValue("KAFKA_BROKER_PORT", "9092")
	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerHost + ":" + brokerPort},
		Topic:   "session",
		GroupID: "dialogs",
		Logger:  l,
	})
	return Consumer{Reader: r, SessionStore: sessionStore}
}

func (c *Consumer) ReadSessionInfo(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		var sesssion models.UserSession
		err = msgpack.Unmarshal(msg.Value, &sesssion)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("received: ", sesssion.UserId)
		_, err = c.SessionStore.CreateSession(&sesssion)
		if err != nil {
			log.Println(err)
		}
	}
}
