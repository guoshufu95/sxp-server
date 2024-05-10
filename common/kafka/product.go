package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	_const "sxp-server/common/const"
	"sxp-server/common/logger"
	"sxp-server/config"
)

type ImplProductHandler struct {
}

func (h *ImplProductHandler) Do(msg kafka.Message) (err error) {
	l := logger.GetLogger()
	l.Info(fmt.Sprintf("receive approve, topic:%s,  partition:%v,   key:%s,  value:%s", msg.Topic, msg.Partition, msg.Key, string(msg.Value)))
	return
}

// NewProductConsumer
//
//	@Description: 开启一个product消费者
//	@param ctx
func NewProductConsumer(ctx context.Context) {
	// 消费者数量，可根据实际情况进行配置
	for i := 0; i < DefaultConsumerNum; i++ {
		go func() {
			m := NewManager(config.Conf.Kafka.Brokers, _const.ProductConsumerTopic, "product-group", config.Conf.Kafka.Async)
			m.Impl = &ImplProductHandler{}
			m.Start(ctx, DefaultConsumerNum)
		}()
	}
}

func NewProductProducer(ctx context.Context, req any, retry int) (err error) {
	w := NewProducer(config.Conf.Kafka.Brokers, config.Conf.Kafka.ProducerTimeOut, _const.TaskProducerTopic, config.Conf.Kafka.Async)
	err = Send2Topic(ctx, w, req, retry)
	return
}
