package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"sxp-server/config"
	"testing"
	"time"
)

func TestReceive(t *testing.T) {
	config.ReadConfig("../../config/sxp.yml")
	wk := NewWrapKafka("", "test", "1")
	type rs struct {
		WrapKafka
		Ch chan kafka.Message
	}
	var tt = rs{
		WrapKafka: wk,
		Ch:        make(chan kafka.Message),
	}
	err := tt.Receive(context.Background(), tt.Ch)
	if err != nil {
		return
	}
	for {
		select {
		case v, _ := <-tt.Ch:
			s := string(v.Value)
			fmt.Println(s)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
