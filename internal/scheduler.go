package internal

import (
	"sync"
	"time"

	"ottplay-foss-server-epg/pkg/helpers"
	"ottplay-foss-server-epg/pkg/mem_storage"
)

const sched_interval = 16 * time.Minute

var (
	sched_ticker *time.Ticker = time.NewTicker(sched_interval)
	_mono_ticker sync.Mutex
	_now_unix    int64
)

func SchedulerCall(store mem_storage.ProvStore) {
	sched_ticker.Reset(sched_interval) // cбрасываем таймер планировщика
	tick_func(time.Now(), store)       // запускаем обновление каналов
}

func tick_func(t time.Time, store mem_storage.ProvStore) {
	_mono_ticker.Lock()
	_now_unix = t.Unix()
	var p *mem_storage.ProviderElement
	for _, p = range store {
		if p.ChTTR <= _now_unix {
			p.ChTTR = _now_unix + int64(p.ChTTL*3600)
			DownloadProvider(store, p)
		}
	}
	_mono_ticker.Unlock()
	helpers.PrintMemUsage("scheduler")
}

func StartSched() {
	// Сброс планировщика
	sched_ticker.Reset(sched_interval)

	for t := range sched_ticker.C {
		tick_func(t, P)
	}
}
