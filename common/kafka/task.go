package kafka

import (
	"context"
	_const "sxp-server/common/const"
	"sxp-server/config"
)

// NewTaskConsumer
//
//	@Description: 开始一个task消费者
func NewTaskConsumer(ctx context.Context) {
	m := NewManager(config.Conf.Kafka.Brokers, _const.TaskConsumerTopic, "task-group", config.Conf.Kafka.Async)
	m.Start(ctx, DefaultConsumerNum)
}

func NewTaskProducer(ctx context.Context, req any, retry int) (err error) {
	w := NewProducer(config.Conf.Kafka.Brokers, config.Conf.Kafka.ProducerTimeOut, _const.TaskProducerTopic, config.Conf.Kafka.Async)
	err = Send2Topic(ctx, w, req, retry)
	return
}
