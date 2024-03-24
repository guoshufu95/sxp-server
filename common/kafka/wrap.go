package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"sxp-server/common/logger"
)

type Processor interface {
	Send(context context.Context, v ...interface{}) error
	Receive(context context.Context) (kafka.Message, error)
}

// WrapKafka
// @Description: 实现Processor接口
type WrapKafka struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
	Log    *logger.ZapLog
}

// NewWrapKafka
//
//	@Description: 返回一个WrapKafka实例
//	@param writeTopic
//	@param readTopic
//	@param groupId
//	@return *WrapKafka
func NewWrapKafka(writeTopic, readTopic, groupId string) WrapKafka {
	manager := NewManager()
	var wk WrapKafka
	if writeTopic != "" {
		wk.Writer = manager.NewProducer(writeTopic)
	}
	if readTopic != "" && groupId != "" {
		wk.Reader = manager.NewConsumer(readTopic, groupId)
	}
	wk.Log = logger.GetLogger()
	return wk
}

// Send
//
//	@Description: 发送
//	@receiver k
//	@param ctx
//	@param data
//	@return err
func (k *WrapKafka) Send(ctx context.Context, data ...interface{}) (err error) {
	//msg := make([]kafka.Message, 0)
	for _, v := range data {
		b, _ := json.Marshal(v)
		err = k.Writer.WriteMessages(ctx, kafka.Message{Value: b})
		if err != nil {
			k.Log.Error("kafka生产者发送消息错误")
			return err
		}
	}
	return
}

// Receive
//
//	@Description: 接收
//	@receiver k
//	@param ctx
//	@param receive
//	@return err
func (k *WrapKafka) Receive(ctx context.Context, receive chan kafka.Message) (err error) {
	go func() {
		var msg kafka.Message
		for {
			msg, err = k.Reader.ReadMessage(ctx)
			if err != nil {
				k.Log.Errorf("kafak消费者错误")
				break
			}
			receive <- msg
		}
		return
	}()
	return
}
