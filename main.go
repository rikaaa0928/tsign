package main

import (
	"flag"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rikaaa0928/tsign/web"
	"github.com/rikaaa0928/tsign/worker"
)

func main() {
	configPath := flag.String("c", "/etc/tsign", "path to cookies")
	port := flag.Int("p", 60080, "port of dashboard")
	flag.Parse()
	timeZone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		timeZone = time.Local
	}
	lastDay := time.Now().In(timeZone).Day()
	w := worker.NewWorker()
	w.AsyncGo()
	um := &atomic.Value{}
	go web.Start(um, *port)
	um.Store(worker.NewUserManager(*configPath))
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
		timeZone, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			timeZone = time.Local
		}
		now := time.Now().In(timeZone)
		if now.Day() != lastDay {
			lastDay = now.Day()
			um.Store(worker.NewUserManager(*configPath))
		}
	}
}
