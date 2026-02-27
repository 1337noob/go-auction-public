package integration_event_bus

import (
	"context"
	"fmt"
	"log"
	"main/pkg"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

const (
	retryDelay = 2 * time.Second
)

type KafkaIntegrationEventBus struct {
	producer           sarama.SyncProducer
	consumer           sarama.ConsumerGroup
	topicPrefix        string
	subscribers        map[EventType][]Handler
	mu                 sync.RWMutex
	marshaller         *EventMarshaller
	transactionManager pkg.TransactionManager
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
}

type KafkaConfig struct {
	Brokers            []string
	TopicPrefix        string
	ConsumerGroup      string
	TransactionManager pkg.TransactionManager
}

func NewKafkaIntegrationEventBus(config KafkaConfig) (*KafkaIntegrationEventBus, error) {
	if config.TransactionManager == nil {
		return nil, fmt.Errorf("config.TransactionManager is required")
	}
	if config.TopicPrefix == "" {
		config.TopicPrefix = "integration-events"
	}
	if config.ConsumerGroup == "" {
		config.ConsumerGroup = "integration-event-consumers"
	}

	// Producer config
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Timeout = 10 * time.Second
	producerConfig.Producer.Idempotent = true
	producerConfig.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(config.Brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Consumer config
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerConfig.Version = sarama.V2_0_0_0

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.ConsumerGroup, consumerConfig)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	bus := &KafkaIntegrationEventBus{
		producer:           producer,
		consumer:           consumer,
		topicPrefix:        config.TopicPrefix,
		subscribers:        make(map[EventType][]Handler),
		marshaller:         NewEventMarshaller(),
		transactionManager: config.TransactionManager,
		ctx:                ctx,
		cancel:             cancel,
	}

	return bus, nil
}

func (b *KafkaIntegrationEventBus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *KafkaIntegrationEventBus) Publish(ctx context.Context, event Event) error {
	topic := b.getTopicName(event.GetType())

	eventData, err := b.marshaller.MarshalEvent(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(string(event.GetType())),
		Value: sarama.ByteEncoder(eventData),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event-type"),
				Value: []byte(string(event.GetType())),
			},
		},
	}

	_, _, err = b.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

func (b *KafkaIntegrationEventBus) getTopicName(eventType EventType) string {
	return fmt.Sprintf("%s-%s", b.topicPrefix, string(eventType))
}

func (b *KafkaIntegrationEventBus) consumeMessages() {
	defer b.wg.Done()

	handler := &kafkaConsumerGroupHandler{
		bus:        b,
		marshaller: b.marshaller,
	}

	for {
		topics := b.getAllTopics()
		if len(topics) == 0 {
			select {
			case <-b.ctx.Done():
				return
			case <-time.After(retryDelay):
				continue
			}
		}

		err := b.consumer.Consume(b.ctx, topics, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
			select {
			case <-b.ctx.Done():
				return
			case <-time.After(retryDelay):
			}
		}
	}
}

func (b *KafkaIntegrationEventBus) getAllTopics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.subscribers))
	for eventType := range b.subscribers {
		topics = append(topics, b.getTopicName(eventType))
	}
	return topics
}

func (b *KafkaIntegrationEventBus) StartConsuming() {
	b.wg.Add(1)
	go b.consumeMessages()
}

func (b *KafkaIntegrationEventBus) Close() error {
	b.cancel()
	b.wg.Wait()

	var errs []error

	if err := b.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}

	if err := b.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing Kafka bus: %v", errs)
	}

	return nil
}

type kafkaConsumerGroupHandler struct {
	bus        *KafkaIntegrationEventBus
	marshaller *EventMarshaller
}

func (h *kafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := h.handleMessage(message); err != nil {
				log.Printf("Failed to handle message from topic %s, partition %d, offset %d: %v",
					message.Topic, message.Partition, message.Offset, err)
				// Не коммитим offset при ошибке
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *kafkaConsumerGroupHandler) handleMessage(message *sarama.ConsumerMessage) error {
	eventType := h.extractEventType(message)
	if eventType == "" {
		return fmt.Errorf("could not determine event type from message in topic %s", message.Topic)
	}

	event, err := h.marshaller.UnmarshalEvent(eventType, message.Value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	h.bus.mu.RLock()
	handlers, exists := h.bus.subscribers[eventType]
	h.bus.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return fmt.Errorf("no handlers found for event type: %s", eventType)
	}

	tx, err := h.bus.transactionManager.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctx := pkg.WithTransaction(tx.Context(), tx)

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction: %v", rollbackErr)
			}
			return fmt.Errorf("handler error: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (h *kafkaConsumerGroupHandler) extractEventType(message *sarama.ConsumerMessage) EventType {
	for _, header := range message.Headers {
		if string(header.Key) == "event-type" {
			return EventType(header.Value)
		}
	}

	// Fallback: extract from topic name
	prefix := h.bus.topicPrefix + "-"
	if len(message.Topic) > len(prefix) {
		return EventType(message.Topic[len(prefix):])
	}

	return ""
}
