package main

import (
	"sync"
	"time"

	"github.com/legolord208/gofilesync/api"
)

var mutexTicker sync.Mutex
var ticker *time.Ticker
var tickerDone chan bool

func scheduleStop() {
	if ticker != nil {
		ticker.Stop()
	}
	if tickerDone != nil {
		tickerDone <- true
	}
}
func schedule(minutes int) {
	scheduleStop()

	if minutes <= 0 {
		return
	}

	mutexTicker.Lock()
	defer mutexTicker.Unlock()

	ticker = time.NewTicker(time.Minute * time.Duration(minutes))
	tickerDone = make(chan bool, 1)

	go func() {
		mutexTicker.Lock()
		defer mutexTicker.Unlock()

		for {
			select {
			case <-ticker.C:
				data.mutexLocations.RLock()
				for _, loc := range data.Locations {
					err := gofilesync.LazySync(loc.Src, loc.Dst)
					if err != nil {
						status(true, err.Error())
					}
				}
				data.mutexLocations.RUnlock()
			case <-tickerDone:
				return
			}
		}
	}()
}
