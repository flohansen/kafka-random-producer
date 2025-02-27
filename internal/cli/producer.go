package cli

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

const (
	KiB = 1024
)

type Flags struct {
	Brokers  string
	Topic    string
	Username string
	Password string
}

type Producer struct {
}

func (p *Producer) Run(ctx context.Context, flags Flags) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": flags.Brokers,
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     flags.Username,
		"sasl.password":     flags.Password,
		"group.id":          "",
	})
	if err != nil {
		return fmt.Errorf("could not create producer: %w", err)
	}

	partitions, err := getNumberPartitions(ctx, producer, flags.Topic)
	if err != nil {
		return fmt.Errorf("could not get number partitions for topic '%s': %w", flags.Topic, err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			key := []byte(uuid.New().String())
			value, err := generateMessage(1 * KiB)
			if err != nil {
				return fmt.Errorf("could not generate random message value: %w", err)
			}

			msg := &kafka.Message{
				Key:   key,
				Value: value,
				TopicPartition: kafka.TopicPartition{
					Topic:     &flags.Topic,
					Partition: shard(key, partitions),
				},
			}

			if err := producer.Produce(msg, nil); err != nil {
				return fmt.Errorf("could not produce message: %w", err)
			}

			log.Printf("produced message, key: %s, len: %d bytes\n", string(key), len(value))
		}
	}
}

func shard(b []byte, n int) int32 {
	h := fnv.New64a()
	h.Write(b)
	return int32(h.Sum64() % uint64(n))
}

func getNumberPartitions(ctx context.Context, producer *kafka.Producer, topic string) (int, error) {
	client, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		return -1, fmt.Errorf("could not create admin client: %w", err)
	}

	dtr, err := client.DescribeTopics(ctx, kafka.NewTopicCollectionOfTopicNames([]string{topic}))
	if err != nil {
		return -1, fmt.Errorf("could not describe topic: %w", err)
	}

	if len(dtr.TopicDescriptions) != 1 {
		return -1, errors.New("no topics in result set")
	}

	n := len(dtr.TopicDescriptions[0].Partitions)
	if n == 0 {
		return -1, errors.New("topic does not exists")
	}

	return n, nil
}

func generateMessage(n int) ([]byte, error) {
	b := make([]byte, n/2)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(b)), nil
}
