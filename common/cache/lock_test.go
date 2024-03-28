package cache

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewLock(t *testing.T) {
	client := IniCache()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				p := recover()
				fmt.Println(p)
			}()
			defer wg.Done()
			token := "token" + strconv.Itoa(i)
			llock := NewLock("lock", token, client, WithMaxIdle(23), WithBlockTimeOut(25), WithExpire(5), WithWatch(true), WithIsBlock(true))
			ctx := context.Background()
			if err := llock.Lock(ctx); err != nil {
				t.Error(err)
				return
			}
			defer func() {
				llock.stopDog()
				err := llock.Unlock(ctx)
				if err != nil {
					return
				}
			}()
			time.Sleep(6 * time.Second) //业务处理
			fmt.Println(0 % 1)          //模拟骚操作
		}(i)
	}
	wg.Wait()
}
