package main

import (
	"time"

	"github.com/legolord208/gofilesync/api"
)

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

	ticker = time.NewTicker(time.Minute * time.Duration(minutes))
	tickerDone = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, loc := range data.Locations {
					err := gofilesync.LazySync(loc.Src, loc.Dst)
					if err != nil {
						status(true, err.Error())
					}
				}
			case <-tickerDone:
				return
			}
		}
	}()
}
