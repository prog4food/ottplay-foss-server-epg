package providers_downloader

import (
	"math/rand"
	"time"

	"github.com/valyala/fasthttp"

	"ottplay-foss-server-epg/libs/config_epg"
	"ottplay-foss-server-epg/libs/helpers"
	ms "ottplay-foss-server-epg/libs/mem_storage"
)

var (
  sched_ticker = time.NewTicker(16 * time.Minute)
  sched_first_run = true
  _now_unix int64
)

func StartJob() {
  // Init
  rwTimeout, _ := time.ParseDuration("10s")
  HttpClient = &fasthttp.Client{
    ReadTimeout:                   rwTimeout,
    WriteTimeout:                  rwTimeout,
    MaxResponseBodySize:           5242880, // 5 MB
    //NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
    DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
    DisablePathNormalizing:        true,
    DialDualStack:                 true,
    // increase DNS cache time to an hour instead of default minute
    Dial: (&fasthttp.TCPDialer{
      Concurrency:      4096,
      DNSCacheDuration: time.Hour,
    }).Dial,
  }
  // Job
  tick_func := func (t time.Time) {
    _now_unix = t.Unix()
    //log.Print("...sheduler run ", _now_unix)
    for i := 0; i < len(config_epg.ConfigData.Providers); i++ {
      ps := &config_epg.ConfigData.Providers[i]
      if (ps.NextUpd >= 0) && (ps.NextUpd <= _now_unix) {
        ps.NextUpd = _now_unix + int64(ps.UpdMins * 60) + int64(rand.Intn(360)) // + Разброс 5 мин
        DownloadProvider(ps)
      }
    }
    helpers.PrintMemUsage("scheduler")
  }
  // First run
  tick_func(time.Now())
  sched_first_run = false
  ms.PO.Sort()
  // Start
  for t := range sched_ticker.C { tick_func(t) }
}
