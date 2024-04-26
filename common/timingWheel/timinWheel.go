package timingWheel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	zaplog "sxp-server/common/logger"
	"sync"
	"time"
)

type ff func(string)

// SxpTask
// @Description: task
type SxpTask struct {
	Key       string        `json:"key"`
	ExecuteAt time.Duration `json:"executeAt"`
	Times     int           `json:"times"`
	Slot      int           `json:"slot"`
	Circle    int           `json:"circle"`
}

// CacheFiled
// @Description: redis的储存字段
//type CacheFiled struct {
//	Key  string   `json:"key"`
//	Task *SxpTask `json:"task"`
//}

// SxpTimingWheel
// @Description: 分布式时间轮struct
type SxpTimingWheel struct {
	log         *zaplog.ZapLog //log
	redisClient *redis.Client  //redis客户端
	interval    time.Duration  //槽位时间间隔
	ticker      *time.Ticker   // 定时器
	currentSlot int            //当前槽位
	slotCount   int            //一轮有多少个槽位
	stopCh      chan struct{}  //停止
	//removeTaskCh chan string    //删除
	addTaskCh chan *SxpTask //添加
	mux       sync.Mutex    //锁
	Run       bool          //运行标志
	fnMap     map[string]ff //任务function,（此处为测试的时候使用，真实场景不建议这么定义）
}

// NewSxpTimeWheel
//
//	@Description: new一个调度器并开启消费协程
//	@param cache
//	@param interval
//	@param slotCount
//	@return *SxpTimingWheel
func NewSxpTimeWheel(ctx context.Context, l *zaplog.ZapLog, client *redis.Client, interval time.Duration, slotCount int) *SxpTimingWheel {
	s := SxpTimingWheel{
		log:         l,
		interval:    interval * time.Second,
		slotCount:   slotCount,
		currentSlot: 0,
		redisClient: client,
		stopCh:      make(chan struct{}),
		//removeTaskCh: make(chan string),
		addTaskCh: make(chan *SxpTask),
		Run:       false,
		fnMap:     make(map[string]ff),
	}
	s.Start(ctx)
	return &s
}

// Start
//
//	@Description: 开始
//	@receiver s
func (s *SxpTimingWheel) Start(ctx context.Context) {
	var once sync.Once
	once.Do(func() { //初始化时间轮盘
		for i := 0; i < s.slotCount; i++ {
			idStr := strconv.Itoa(i) + ":timingWheel"
			s.redisClient.HSet(ctx, idStr, "", "")
		}
		s.ticker = time.NewTicker(s.interval)
		s.mux.Lock()
		defer s.mux.Unlock()
		s.Run = true
		go s.handle(ctx)
	})
}

// handle
//
//	@Description: 分发逻辑
//	@receiver s
func (s *SxpTimingWheel) handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // 停止
			s.ticker.Stop()
			s.log.Info("##################### 任务队列退出 #####################")
			return
		case task := <-s.addTaskCh: //新增task
			s.addTask(task)
		//case key := <-s.removeTaskCh: //删除task
		//	s.removeTask(key)
		case <-s.ticker.C: //定时消费
			s.execute()
		}
	}
}

// CreateTask
//
//	@Description: 创建task
//	@receiver s
//	@param key
//	@param job
//	@param executeAt
//	@param times
//	@return err
func (s *SxpTimingWheel) CreateTask(key string, job func(string), executeAt time.Duration, times int) (task *SxpTask, err error) {
	if key == "" {
		return nil, errors.New("key不能为空")
	}
	if int64(executeAt) < time.Now().Unix() {
		return nil, errors.New("executeAt时间错误")
	}
	key = fmt.Sprintf("%s:%d", key, executeAt)
	task = &SxpTask{
		Key:       key,
		Times:     times,
		ExecuteAt: executeAt * time.Second,
	}
	s.fnMap[key] = job
	s.addTaskCh <- task
	return
}

// addTask
//
//	@Description: 添加task的具体逻辑
//	@receiver s
//	@param task
func (s *SxpTimingWheel) addTask(task *SxpTask) {
	ctx := context.Background()
	slot, circle := s.cal(task.ExecuteAt)
	task.Slot = slot
	task.Circle = circle
	key := strconv.Itoa(slot) + ":timingWheel"
	b, err := json.Marshal(task)
	if err != nil {
		s.log.Errorf("序列化错误:%s", err.Error())
		return
	}
	err = s.redisClient.HSet(ctx, key, task.Key, b).Err()
	if err != nil {
		s.log.Errorf("redis添加失败: %s", err.Error())
		return
	}
	s.log.Info("新增task成功!")
	return
}

//// removeTask
////
////	@Description: 删除一个task
////	@receiver s
////	@param key
//func (s *SxpTimingWheel) removeTask(keys string) {
//	kl := strings.Split(keys, "-")
//	err := s.redisClient.HDel(context.Background(), kl[0], kl[1]).Err()
//	if err != nil {
//		s.log.Errorf("删除键值失败:%s", err.Error())
//		return
//	}
//}

// execute
//
//	@Description: 执行任务
//	@receiver s
func (s *SxpTimingWheel) execute() {
	key := s.buildKey()
	m := s.redisClient.HGetAll(context.Background(), key).Val()
	if m != nil {
		for k, v := range m {
			if k == "" || v == "" {
				continue
			}
			var task SxpTask
			err := json.Unmarshal([]byte(v), &task)
			if err != nil {
				s.log.Errorf("json反序列化失败:%s", err.Error())
				continue
			}
			if task.Circle > 0 {
				task.Circle--
				by, er := json.Marshal(task)
				if er != nil {
					s.log.Errorf("json序列化失败: %s", err.Error())
					continue
				}
				_, err = s.redisClient.HSet(context.Background(), key, k, by).Result()
				if err != nil {
					s.log.Errorf("redis hset错误:%s", er.Error())
				}
				continue
			}
			//todo 更新数据库或做超时控制等
			s.fnMap[task.Key]("测试task" + v)
			delete(s.fnMap, k)
			err = s.redisClient.HDel(context.Background(), key, k).Err()
			if err != nil {
				s.log.Errorf("删除键值失败:%s", err.Error())
				return
			}
		}
	}
	s.incrCurrentSlot()
}

// cal
//
//	@Description: 计算当前task的slot和圈数
//	@receiver s
//	@param executeAt
//	@return slot
//	@return circle
func (s *SxpTimingWheel) cal(executeAt time.Duration) (slot, circle int) {
	delay := int(int64(executeAt.Seconds()) - time.Now().Unix())
	circleTime := s.slotCount * int(s.interval.Seconds())
	circle = delay / circleTime
	steps := delay / int(s.interval.Seconds())
	slot = (s.currentSlot + steps) % s.slotCount
	return
}

// incrCurrentSlot
//
//	@Description: 槽位自增
//	@receiver s
func (s *SxpTimingWheel) incrCurrentSlot() {
	s.currentSlot = (s.currentSlot + 1) % s.slotCount
}

func (s *SxpTimingWheel) buildKey() string {
	key := strconv.Itoa(s.currentSlot) + ":timingWheel"
	return key
}
