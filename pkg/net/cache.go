package net

import (
	"fmt"
	"time"
)

type cache map[string]chan interface{}
type sdata map[string]interface{}
type handler func(sdata, cache_config) (interface{}, error)
type task map[string]Work

type cache_config struct {
	CacheLifetime int
	QueryDelay    int
	QueryRetries  int
}

type Work struct {
	Data    sdata
	Channel chan interface{}
	Handler handler
}

func WorkerPool(conf cache_config, work chan task, workcache cache, shutdown chan int) {
	for {
		//check if shutting down
		select {
		case nexttask := <-work:
			//perform work
			for hash, curtask := range nexttask {
				if _, ok := workcache[hash]; !ok { //cache miss
					workcache[hash] = curtask.Channel
					go WorkerThread(conf, curtask.Handler, workcache[hash], curtask.Data, shutdown)
				} else { //cache hit
					select {
					case res := <-workcache[hash]: //cache is still current
						curtask.Channel <- res
					default: //cache timed out
						go WorkerThread(conf, curtask.Handler, workcache[hash], curtask.Data, shutdown)
						curtask.Channel <- <-workcache[hash]
					}
					continue
				}
			}
		case <-shutdown:
			shutdown <- 1
			return
		}
	}
}

func WorkerThread(conf cache_config, hfunc handler, entry chan interface{}, data interface{}, shutdown chan int) {
	//request result until no error is returned
	timeout := conf.CacheLifetime

	res, err := hfunc(data.(sdata), conf)
	retry := 1

	for err != nil {
		time.Sleep(time.Duration(retry*conf.QueryDelay) * time.Millisecond)
		res, err = hfunc(data.(sdata), conf)
		retry++
		if err != nil && retry > conf.QueryRetries {
			timeout = retry * conf.QueryDelay
			res = fmt.Sprintf("%+v", err)
			err = nil
		}
	}

	//continue to serve data from cache for lifetime
	lifetime := time.NewTimer(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-shutdown:
			shutdown <- 1
			return
		case <-lifetime.C:
			return
		case entry <- res:
		}
	}
}
