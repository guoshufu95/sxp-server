package timingWheel

import (
	"context"
	"github.com/redis/go-redis/v9"
	"net/http"
	"sync"
	"time"
)

type SxpTask struct {
	Key string `json:"key"`
}

// SxpTimeWheel
//
//	@Description: redis版本时间轮struct
type SxpTimeWheel struct {
	redisClient *redis.Client
	httpClient  *http.Client
	stopCh      chan struct{}
	ticker      *time.Ticker
}

// NewSxpTimeWheel
//
//	@Description: 初始化
//	@param cache
//	@param client
//	@return *SxpTimeWheel
func NewSxpTimeWheel(cache *redis.Client, client *http.Client) *SxpTimeWheel {
	s := SxpTimeWheel{
		ticker:      time.NewTicker(time.Second),
		redisClient: cache,
		httpClient:  client,
		stopCh:      make(chan struct{}),
	}
	return &s
}

// Run
//
//	@Description: 监听
//	@receiver s
func (s *SxpTimeWheel) Run() {
	for {
		select {
		case <-s.stopCh:
			return
		case <-s.ticker.C:

		}
	}
}

// Stop
//
//	@Description: 停止
//	@receiver s
func (s *SxpTimeWheel) Stop() {
	var once *sync.Once
	once.Do(func() {
		close(s.stopCh)
		s.ticker.Stop()
	})
}

// AddTask
//
//	@Description: 新增任务
//	@receiver s
//	@param ctx
//	@param key
//	@param task
//	@param t
func (s *SxpTimeWheel) AddTask(ctx context.Context, key string, task *SxpTask, t time.Time) {

}
