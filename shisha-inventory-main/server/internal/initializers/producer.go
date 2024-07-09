package initializers

import (
	"context"
	"encoding/json"
	"server/internal/models"

	zlog "github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{client: client, topic: topic}, nil
}
func (p *Producer) SendTransferMessage(ctx context.Context, user, target string, amount int) {
	msg := models.TranfserMessage{User: user, Type: "transfer", Target: target, Amount: amount}
	b, _ := json.Marshal(msg)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, func(_ *kgo.Record, err error) {
		if err != nil {
			zlog.Printf("record had a produce error: %v\n", err)
		}
	})
}

func (p *Producer) SendUploadMessage(ctx context.Context, user, image_uuid string) {
	msg := models.UploadMessage{User: user, Type: "upload", Image_uuid: image_uuid}
	b, _ := json.Marshal(msg)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, func(_ *kgo.Record, err error) {
		if err != nil {
			zlog.Printf("record had a produce error: %v\n", err)
		}
	})
}

func (p *Producer) SendBuyMessage(ctx context.Context, user, image_uuid string, amount int) {
	msg := models.BuyMessage{User: user, Type: "buy", Image_uuid: image_uuid, Amount: amount}
	b, _ := json.Marshal(msg)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, func(_ *kgo.Record, err error) {
		if err != nil {
			zlog.Printf("record had a produce error: %v\n", err)
		}
	})
}

func (p *Producer) Close() {
	p.client.Close()
}

func RedPandaReady(ctx context.Context, brokers []string) bool {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		zlog.Printf("Error creating client: %v", err)
		zlog.Fatal()
	}
	err = client.Ping(ctx)
	if err != nil {
		zlog.Printf("Error ping RedPandas client: %v", err)
		zlog.Fatal()
		return false
	} else {
		return true
	}
}
