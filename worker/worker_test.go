package worker

import (
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := &worker{c: make(chan *workerData)}
	w.AsyncGo()
	u := NewUserManager("../test")
	wg := sync.WaitGroup{}
	for _, v := range u.UserMap() {
		wg.Add(1)
		dm := v.DataMap()
		a := v.a
		go func() {
			for _, d := range dm {
				w.AsyncNotify(&workerData{d: d, act: a, f: func(showData ShowData) {
					t.Log(showData)
				}})
				time.Sleep(time.Second)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second)
}