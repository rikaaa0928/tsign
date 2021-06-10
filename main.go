package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/rikaaa0928/tsign/worker"
)

func main() {
	lastDay := time.Now().Day()
	w := worker.NewWorker()
	w.AsyncGo()
	um := atomic.Value{}
	//u := worker.NewUserManager("../test")
	um.Store(worker.NewUserManager("./test"))
	go func() {
		for {
			u := um.Load().(*worker.UserMgr)
			wg := sync.WaitGroup{}
			wg.Add(1)
			var once sync.Once
			for _, v := range u.UserMap() {
				v := v
				go func() {
					v.FeedWorker(w)
					time.Sleep(time.Minute * 30)
					once.Do(func() {
						wg.Done()
					})
				}()
			}
			wg.Wait()
		}
	}()

	ticker := time.Tick(time.Second)
	for {
		<-ticker
		now := time.Now()
		if now.Day() != lastDay {
			lastDay = now.Day()
			um.Store(worker.NewUserManager("./test"))
		}
	}
}
