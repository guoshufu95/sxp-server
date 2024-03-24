package kafka

import (
	"github.com/segmentio/kafka-go"
	"sxp-server/config"
	"time"
)

type Manager struct {
	Brokers []string
}

type MQ struct {
	Topic    string
	GroupId  string
	Producer *kafka.Writer
	Consumer *kafka.Reader
}

// NewManager
//
//	@Description: 创建一个manager实例
//	@return Manager
func NewManager() *Manager {
	var m Manager
	m.Brokers = config.Conf.Kafka.Brokers
	return &m
}

// NewProducer
//
//	@Description: 新建一个生产者
//	@receiver m
//	@param topic
//	@return *kafka.Writer
func (m *Manager) NewProducer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(m.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           time.Duration(config.Conf.Kafka.ProducerTimeOut),
		RequiredAcks:           kafka.RequiredAcks(config.Conf.Kafka.Ack),
		AllowAutoTopicCreation: true,
	}
}

// NewConsumer
//
//	@Description: 新建一个消费者
//	@receiver m
//	@param topic
//	@param groupId
//	@return *kafka.Reader
func (m *Manager) NewConsumer(topic, groupId string) *kafka.Reader {
	// TODO reader 优雅关闭
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          m.Brokers,
		Topic:            topic,
		GroupID:          groupId,
		ReadBatchTimeout: time.Duration(config.Conf.Kafka.ConsumerTimeOut),
		StartOffset:      kafka.FirstOffset,
	})
	return reader
}
