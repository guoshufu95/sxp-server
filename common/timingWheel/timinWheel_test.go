package timingWheel

import (
	"context"
	"fmt"
	"sxp-server/common/initial"
	"testing"
	"time"
)

func TestNewSxpTimeWheel(t *testing.T) {
	app := initial.App
	fn := func(s string) {
		fmt.Println("这是测试数据： ", s)
		return
	}
	ctx := context.Background()
	s := NewSxpTimeWheel(ctx, app.Logger, app.Cache, 1, 5)
	_, _ = s.CreateTask("testfn", fn, time.Duration(1714206174), 1)
	for {
		time.Sleep(1 * time.Second)
	}
}

//func TestSxpTimingWheel_removeTask(t *testing.T) {
//	fn := func(s string) {
//		fmt.Println("这是测试数据： ", s)
//		return
//	}
//	app := initial.App
//	ctx := context.Background()
//	s := NewSxpTimeWheel(ctx, app.Logger, app.Cache, 1, 5)
//	task, _ := s.CreateTask("testfn", fn, time.Duration(1712643328), 1)
//	key0 := fmt.Sprintf("%s:%d", "testfn", 1712643328)
//	key := fmt.Sprintf("%s-%s", strconv.Itoa(task.Slot)+":timingWheel", key0)
//	go func() {
//		s.removeTaskCh <- key
//	}()
//	time.Sleep(30 * time.Second)
//	fmt.Println("### 删除任务 ###")
//}
