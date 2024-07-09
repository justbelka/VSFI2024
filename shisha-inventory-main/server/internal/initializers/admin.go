package initializers

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Admin struct {
	client *kadm.Client
}

func NewAdmin(brokers []string) (*Admin, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	admin := kadm.NewClient(client)
	return &Admin{client: admin}, nil
}
func (a *Admin) TopicExists(ctx context.Context, topic string) (bool, error) {
	topicsMetadata, err := a.client.ListTopics(ctx)
	if err != nil {
		return false, err
	}
	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true, nil
		}
	}
	return false, nil
}
func (a *Admin) CreateTopic(ctx context.Context, topic string) error {
	resp, err := a.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		return err
	}
	for _, ctr := range resp {
		if ctr.Err != nil {
			return err
		}

		log.Info().Str("topic", ctr.Topic).Msg("Created topic")
	}
	return nil
}
func (a *Admin) Close() {
	a.client.Close()
}
