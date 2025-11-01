package mq

import (
	"devmgr/internal/consts"
	"devmgr/internal/service"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func New(brokers string, groupId string) service.IMQClient {
	return &sMQClient{
		pConfig: &kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"acks":              "all",
			"retries":           3,
		},
		cConfig: &kafka.ConfigMap{
			"bootstrap.servers":                  brokers,
			"group.id":                           groupId, // 消费者组标识
			"auto.offset.reset":                  "latest",
			"enable.auto.commit":                 true,
			"auto.commit.interval.ms":            1000,
			"socket.timeout.ms":                  5000,
			"socket.connection.setup.timeout.ms": 5000,
		},
		subscriptions:     make(map[string]consts.MessageHandler),
		running:           true,
		reconnectInterval: 5 * time.Second,
	}
}

type Publisher interface {
	Publish(topic string, key []byte, payload []byte) error
}

type sMQClient struct {
	pConfig           *kafka.ConfigMap
	cConfig           *kafka.ConfigMap
	producer          *kafka.Producer
	consumer          *kafka.Consumer
	subscriptions     map[string]consts.MessageHandler // Topic to handler mapping
	running           bool                             // Flag to control running state
	mu                sync.RWMutex                     // Protects consumer and producer pointers
	reconnectInterval time.Duration                    // Interval for reconnecting
}

// Init initializes the producer and starts the consume loop.
func (k *sMQClient) Start() error {
	producer, err := kafka.NewProducer(k.pConfig)
	if err != nil {
		return err
	}
	k.producer = producer

	// Start the auto-reconnecting consume loop.
	go k.startConsumeLoop()
	log.Printf("Create Kafka connection done")
	return nil
}

// startConsumeLoop runs a loop that automatically reconnects on failure.
func (k *sMQClient) startConsumeLoop() {
	for k.isRunning() {
		err := k.connectAndConsume()
		if err != nil {
			log.Printf("Kafka connection/consume failed: %v. Retrying in %v...", err, k.reconnectInterval)
		}
		if k.isRunning() {
			time.Sleep(k.reconnectInterval)
		}
	}
}

// connectAndConsume creates a new consumer, subscribes to topics, and polls for messages.
func (k *sMQClient) connectAndConsume() error {
	consumer, err := kafka.NewConsumer(k.cConfig)
	if err != nil {
		return err
	}

	// Set the new consumer.
	k.mu.Lock()
	k.consumer = consumer
	k.mu.Unlock()

	// Extract topic list from the handler map.
	topics := make([]string, 0, len(k.subscriptions))
	for topic := range k.subscriptions {
		topics = append(topics, topic)
	}

	//订阅Kafka主题
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		consumer.Close()
		return err
	}

	// 开始消费消息
	run := true
	for run && k.isRunning() {

		// 轮询Kafka消息
		ev := k.consumer.Poll(500)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			// Dispatch message to its registered handler.
			if handler, exists := k.subscriptions[*e.TopicPartition.Topic]; exists && handler != nil {
				handler(*e.TopicPartition.Topic, e.Key, e.Value)
			} else {
				log.Printf("[Consumed] No handler for topic '%s': Key=%s, Value=%s",
					*e.TopicPartition.Topic, string(e.Key), string(e.Value))
			}

		case kafka.Error:
			log.Printf("Kafka error: %v", e)
			if e.IsFatal() {
				run = false // Exit to trigger reconnection
			}
		}
	}

	// Clean up the consumer.
	k.mu.Lock()
	if k.consumer == consumer {
		k.consumer.Close()
		k.consumer = nil
	}
	k.mu.Unlock()

	return nil
}

// Publish sends a message synchronously.
func (k *sMQClient) Publish(topic string, key []byte, payload []byte) error {
	k.mu.RLock()
	producer := k.producer
	k.mu.RUnlock()

	if producer == nil {
		return kafka.NewError(kafka.ErrQueueFull, "producer not ready", false)
	}

	// 生产Kafka消息
	deliveryChan := make(chan kafka.Event, 1)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: payload,
	}, deliveryChan)

	if err != nil {
		close(deliveryChan)
		return err
	}

	//等待交付报告或错误
	e := <-deliveryChan
	m := e.(*kafka.Message)
	close(deliveryChan)

	// 检查交付错误，即生成者方确保发送到Broker的消息不丢失
	// 但可能重复发送，如实际发成功了，但是Broker的ACK返回给生产者时出现了网络Error
	// 从而重试后，导致消息重复发送，这时候需要下游做好幂等处理
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

// Register subscription, execute all registered subscriptions after connecting to Broker.
func (k *sMQClient) RegisterSubscription(topic string, handler consts.MessageHandler) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	k.subscriptions[topic] = handler
}

// Close shuts down the producer and consumer gracefully.
func (k *sMQClient) Close() {
	k.mu.Lock()
	k.running = false
	k.mu.Unlock()

	if k.producer != nil {
		k.producer.Close()
	}
	log.Println("Kafka client closed")
}

// isRunning checks if the client is still running.
func (k *sMQClient) isRunning() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.running
}
