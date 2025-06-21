package kafkaclient

import (
	"context"
	"fmt"
	"github.com/Sanchir01/currency-wallet/pkg/utils"

	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	producer *kafka.Writer
	messages chan kafka.Message
	done     chan struct{}
}

func ensureTopicExists(brokers []string, topic string, partitions int, replicationFactor int, ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka broker: %w", err)
	}
	defer conn.Close()
	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		conn.SetWriteDeadline(deadline)
	} else {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:         topic,
			NumPartitions: partitions,
			
			ReplicationFactor: replicationFactor,
		},
	}

	return conn.CreateTopics(topicConfigs...)
}

func NewProducer(brokers []string, topic string, retries int, ctx context.Context) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers provided")
	}

	err := utils.DoWithTries(func() error {
		var err error
		err = ensureTopicExists(brokers, topic, 1, 1, ctx)
		if err != nil {
			return err
		}
		return nil
	}, retries, 5*time.Second)

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	p := &Producer{
		producer: writer,
		messages: make(chan kafka.Message, 100),
		done:     make(chan struct{}),
	}

	go p.run()

	return p, nil
}
func (p *Producer) Produce(message string, value []byte) error {
	select {
	case p.messages <- kafka.Message{Value: value, Key: []byte(message)}:
		return nil
	default:
		return fmt.Errorf("канал сообщений переполнен")
	}
}
func (p *Producer) run() {
	for {
		select {
		case msg := <-p.messages:
			err := p.producer.WriteMessages(context.Background(), msg)
			if err != nil {
				log.Printf("❌ Kafka write error: %v", err)
			} else {
				log.Printf("✅ Kafka message sent: key=%s value=%s", string(msg.Key), string(msg.Value))
			}
		case <-p.done:
			return
		}
	}
}
func (p *Producer) Close() error {
	defer close(p.done)
	defer close(p.messages)
	return p.producer.Close()
}
