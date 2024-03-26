package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
	_const "sxp-server/common/const"
	"sxp-server/common/logger"
	"sxp-server/config"
	"sync"
	"time"
)

// init
//
//	@Description: 创建一个消费者就新增一个kv
func init() {
	Cm.Filed[_const.TaskConsumerTopic] = NewTaskConsumer       //task消费者
	Cm.Filed[_const.ProductConsumerTopic] = NewProductConsumer //product消费者
}

// Manager
// @Description: 消费者manager
type Manager struct {
	Ch      chan interface{}
	Log     *logger.ZapLog
	Brokers []string
	Group   string
	Topic   string
	//Manual  bool
	Reader *kafka.Reader
}

// ConsumerMap
// @Description: 全局topic-func map
type ConsumerMap struct {
	Filed map[string]func(ctx context.Context)
	Lock  sync.Mutex
}

var (
	Cm = ConsumerMap{
		Filed: make(map[string]func(ctx context.Context)),
	}
	DefaultConsumerNum = 4 //消费者组开启消费者默认个数
	Once               sync.Once
)

// NewManager
//
//	@Description: 创建一个manager实例
//	@return Manager
func NewManager(brokers []string, topic, group string, async bool) *Manager {
	m := &Manager{}
	m.Ch = make(chan interface{})
	m.Log = logger.GetLogger()
	m.Brokers = brokers
	m.Group = group
	m.Topic = topic
	m.Reader = NewConsumer(brokers, topic, group)
	return m
}

// Start
//
//	@Description: 开启消费
//	@receiver m
//	@param ctx
//	@param wg
//	@param worker
func (m *Manager) Start(ctx context.Context, worker int) {
	for i := 0; i < worker; i++ {
		go m.Consume(ctx)
	}
}

// Consume
//
//	@Description: 监听消费数据，reader优雅关闭
//	@receiver m
//	@param ctx
func (m *Manager) Consume(ctx context.Context) {
	go m.receive(ctx, m.Reader)
	for {
		select {
		case <-ctx.Done():
			Once.Do(printStop)
			_ = m.Reader.Close()
			return
		case val := <-m.Ch:
			e, ok := val.(error)
			if ok {
				m.Log.Errorf("消费者错误: %s", e.Error())
				_ = m.Reader.Close()
				return
			}
			vv, ok := val.(kafka.Message)
			if ok {
				//todo 加入自己的业务逻辑
				m.Log.Info(string(vv.Value))
			}
		}
	}
}

// receive
//
//	@Description: 自由实现消费逻辑
//	@receiver m
//	@param ctx
//	@param reader
func (m *Manager) receive(ctx context.Context, reader *kafka.Reader) {
	for {
		value, err := reader.ReadMessage(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "context canceled") {
				m.Ch <- err
			}
			break
		}
		m.Ch <- value
	}
	return
}

// Send2Topic
//
//	@Description: 生产者发送数据
//	@param ctx
//	@param w
//	@param req
//	@param retry
//	@return err
func Send2Topic(ctx context.Context, w *kafka.Writer, req any, retry int) (err error) {
	l := logger.GetLogger()
	for i := 0; i < retry; i++ {
		by, _ := json.Marshal(req)
		err = w.WriteMessages(ctx, kafka.Message{
			Value: by,
		})
		if err != nil {
			l.Info("重试中: %d", i+1)
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
	l.Errorf("%s 生产者发送失败：%s", w.Topic, err.Error())
	return
}

// NewProducer
//
//	@Description: 生产者初始化
//	@param brokers
//	@param timeout
//	@param topic
//	@param async
//	@return *kafka.Writer
func NewProducer(brokers []string, timeout int, topic string, async bool) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           time.Duration(timeout),
		RequiredAcks:           kafka.RequiredAcks(config.Conf.Kafka.Ack),
		AllowAutoTopicCreation: true,
		Async:                  async,
	}
}

// NewConsumer
//
//	@Description: 消费者初始化
//	@param brokers
//	@param topic
//	@param groupId
//	@return *kafka.Reader
func NewConsumer(brokers []string, topic, groupId string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          brokers,
		Topic:            topic,
		GroupID:          groupId,
		ReadBatchTimeout: time.Duration(config.Conf.Kafka.ConsumerTimeOut),
		StartOffset:      kafka.LastOffset,
	})
	return reader
}

// StartKafkaConsume
//
//	@Description: 启动所有消费者
func StartKafkaConsume(ctx context.Context) {
	log := logger.GetLogger()
	for topic, f := range Cm.Filed {
		log.Info(fmt.Sprintf("%s 消费者启动", topic))
		f(ctx)
	}
}

// printStop
//
//	@Description: 打印退出消日志
func printStop() {
	l := logger.GetLogger()
	l.Info("kafka消费者退出")
}
