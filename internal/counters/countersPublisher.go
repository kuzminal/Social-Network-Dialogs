package counters

import (
	"Social-Net-Dialogs/internal/helper"
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Publisher struct {
	Writer *kafka.Writer
}

func NewCountersPublisher() Publisher {
	brokerHost := helper.GetEnvValue("KAFKA_BROKER_HOST", "localhost")
	brokerPort := helper.GetEnvValue("KAFKA_BROKER_PORT", "9092")
	l := log.New(os.Stdout, "kafka writer: ", 0)
	w := kafka.Writer{
		Addr:                   kafka.TCP(brokerHost + ":" + brokerPort),
		Topic:                  "counters",
		Logger:                 l,
		AllowAutoTopicCreation: true,
	}
	return Publisher{Writer: &w}
}

func (s *Publisher) SendMessageInfo(ctx context.Context, message []byte) {
	var err error
	const retries = 3
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for i := 0; i < retries; i++ {
		err = s.Writer.WriteMessages(ctx, kafka.Message{
			//Key: []byte(strconv.Itoa(i)),
			Value: message,
		})
		if errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Printf("unexpected error %v", err)
		}
		break
	}
}
