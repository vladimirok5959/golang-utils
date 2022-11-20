package servauth

import (
	"sync"
	"time"
)

type Requests struct {
	sync.RWMutex
	counter   map[string]int
	lastTime  map[string]int64
	cleanTime int64
}

func (r *Requests) Count(ip string) int {
	r.Lock()
	defer r.Unlock()
	if v, ok := r.counter[ip]; ok {
		return v
	}
	return 0
}

func (r *Requests) SetCount(ip string, count int) {
	r.Lock()
	defer r.Unlock()
	r.counter[ip] = count
}

func (r *Requests) Time(ip string) int64 {
	r.Lock()
	defer r.Unlock()
	if v, ok := r.lastTime[ip]; ok {
		return v
	}
	return 0
}

func (r *Requests) SetTime(ip string, time int64) {
	r.Lock()
	defer r.Unlock()
	r.lastTime[ip] = time
}

func (r *Requests) Cleanup() {
	r.Lock()
	defer r.Unlock()
	r.counter = map[string]int{}
	r.lastTime = map[string]int64{}
	r.cleanTime = time.Now().UTC().Unix()
}
