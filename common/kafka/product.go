package kafka

import (
	"context"
	_const "sxp-server/common/const"
	"sxp-server/config"
)

// NewProductConsumer
//
//	@Description: 开启一个product消费者
//	@param ctx
func NewProductConsumer(ctx context.Context) {
	m := NewManager(config.Conf.Kafka.Brokers, _const.ProductConsumerTopic, "product-group", config.Conf.Kafka.Async)
	m.Start(ctx, DefaultConsumerNum)
}
