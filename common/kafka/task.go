package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	_const "sxp-server/common/const"
	"sxp-server/common/logger"
	"sxp-server/config"
)

type ImplTaskHandler struct {
}

func (h *ImplTaskHandler) Do(msg kafka.Message) (err error) {
	l := logger.GetLogger()
	l.Info(fmt.Sprintf("receive approve, topic:%s,  partition:%v,   key:%s,  value:%s", msg.Topic, msg.Partition, msg.Key, string(msg.Value)))
	return
}

// NewTaskConsumer
//
//	@Description: 开始一个task消费者
func NewTaskConsumer(ctx context.Context) {
	// 消费者数量，可根据实际情况进行配置
	for i := 0; i < DefaultConsumerNum; i++ {
		go func() {
			m := NewManager(config.Conf.Kafka.Brokers, _const.TaskConsumerTopic, "task-group", config.Conf.Kafka.Async)
			m.Impl = &ImplTaskHandler{}
			m.Start(ctx, DefaultConsumerNum)
		}()
	}
}

func NewTaskProducer(ctx context.Context, req any, retry int) (err error) {
	w := NewProducer(config.Conf.Kafka.Brokers, config.Conf.Kafka.ProducerTimeOut, _const.TaskProducerTopic, config.Conf.Kafka.Async)
	err = Send2Topic(ctx, w, req, retry)
	return
}
