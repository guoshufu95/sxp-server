package timingWheel

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	defaultInterval  = 1  //默认每个槽位的时间间隔
	defaultSlotCount = 60 //默认多少个槽位
)

// Task
// @Description: 单个任务
type Task struct {
	key       string
	job       func(string)
	executeAt time.Duration
	times     int
	slot      int
	circle    int
}

// SingleSXPTimeWheel
// @Description: 时间轮结构
type SingleSXPTimeWheel struct {
	interval     time.Duration //槽位时间间隔
	ticker       *time.Ticker  // 定时器
	currentSlot  int           //当前槽位
	slotCount    int           //一轮有多少个槽位
	slots        []*list.List  //双向链表,存入redis
	stopCh       chan struct{} //停止
	removeTaskCh chan string   //删除
	addTaskCh    chan *Task    //添加
	taskRecords  sync.Map      //key->task
	mux          sync.Mutex    //锁
	Run          bool          //运行标志
}

// NewSingleSXPTimeWheel
//
//	@Description: 初始化
//	@param interval
//	@param slotCount
func NewSingleSXPTimeWheel(interval time.Duration, slotCount int) {
	t := &SingleSXPTimeWheel{
		interval:     interval,
		currentSlot:  0,
		slotCount:    slotCount,
		slots:        make([]*list.List, slotCount),
		stopCh:       make(chan struct{}),
		removeTaskCh: make(chan string),
		addTaskCh:    make(chan *Task),
		Run:          false,
	}
	if interval <= 0 {
		t.interval = time.Duration(defaultInterval)
	}
	if t.slotCount <= 0 {
		t.slotCount = defaultSlotCount
	}
}

// Start
//
//	@Description: 开始任务
//	@receiver s
func (s *SingleSXPTimeWheel) Start() {
	var once *sync.Once
	once.Do(func() { //程序运行期间只需要初始化一次
		for i := 0; i < s.slotCount; i++ {
			s.slots[i] = list.New()
		}
		s.ticker = time.NewTicker(s.interval)
		s.mux.Lock()
		defer s.mux.Unlock()
		s.Run = true
		go s.handle()
	})
}

// handle
//
//	@Description: handle
//	@receiver s
func (s *SingleSXPTimeWheel) handle() {
	for {
		select {
		case <-s.stopCh:
			return
		case task := <-s.addTaskCh:
			s.addTask(task)
		case key := <-s.removeTaskCh:
			s.removeTask(key)
		case <-s.ticker.C:
			s.execute()
		}
	}
}

// AddTask
//
//	@Description: 添加task
//	@receiver s
//	@param key
//	@param job
//	@param executeAt
//	@param times
//	@return error
func (s *SingleSXPTimeWheel) AddTask(key string, job func(string), executeAt time.Duration, times int) error {
	if key == "" {
		return errors.New("key is empty")
	}
	if executeAt < s.interval {
		return errors.New("key is empty")
	}
	_, ok := s.taskRecords.Load(key)
	if ok {
		return errors.New("key of job already exists")
	}
	task := &Task{
		key:       key,
		job:       job,
		times:     times,
		executeAt: executeAt,
	}
	s.addTaskCh <- task
	return nil
}

// addTask
//
//	@Description:
//	@receiver s
//	@param task
func (s *SingleSXPTimeWheel) addTask(task *Task) {
	slot, circle := s.cal(task.executeAt)
	task.slot = slot
	task.circle = circle
	ele := s.slots[slot].PushBack(task)
	s.taskRecords.Store(task.key, ele)
}

// removeTask
//
//	@Description: 移除任务
//	@receiver s
//	@param key
func (s *SingleSXPTimeWheel) removeTask(key string) {
	taskRec, ok := s.taskRecords.Load(key)
	if !ok {
		return
	}
	ele := taskRec.(*list.Element)
	task, _ := ele.Value.(*Task)
	s.slots[task.slot].Remove(ele)
	s.taskRecords.Delete(key)
}

// execute
//
//	@Description: 执行任务
//	@receiver s
func (s *SingleSXPTimeWheel) execute() {
	taskList := s.slots[s.currentSlot]
	if taskList != nil {
		for ele := taskList.Front(); ele != nil; {
			taskEle, _ := ele.Value.(*Task)
			if taskEle.circle > 0 {
				taskEle.circle--
				ele = ele.Next()
				continue
			}
			go taskEle.job(taskEle.key)
			s.taskRecords.Delete(taskEle.key)
			taskList.Remove(ele)

			if taskEle.times-1 > 0 {
				taskEle.times--
				s.addTask(taskEle)
			}
			if taskEle.times == -1 {
			}
			ele = ele.Next()
		}
	}
	s.incrCurrentSlot()
}

func (s *SingleSXPTimeWheel) incrCurrentSlot() {
	s.currentSlot = (s.currentSlot + 1) % len(s.slots)
}

// cal
//
//	@Description: 计算新增task的槽位和圈数
//	@receiver t
//	@param executeAt
//	@return slot
//	@return circle
func (t *SingleSXPTimeWheel) cal(executeAt time.Duration) (slot, circle int) {
	delay := int(executeAt.Seconds())
	circleTime := len(t.slots) * int(t.interval.Seconds())
	circle = delay / circleTime
	steps := delay / int(t.interval.Seconds())
	slot = (t.currentSlot + steps) % len(t.slots)
	return
}

// Stop
//
//	@Description: 停止
//	@receiver s
func (s *SingleSXPTimeWheel) Stop() {
	if s.Run {
		s.mux.Lock()
		s.Run = false
		s.ticker.Stop()
		s.mux.Unlock()
		close(s.stopCh)
	}
}
