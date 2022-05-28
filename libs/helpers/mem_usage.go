package helpers

import (
	"runtime"

	"github.com/rs/zerolog/log"
)

func ForceGC() { runtime.GC() }

func PrintMemUsage(s string) {
  var m runtime.MemStats
  bToMb := func (b uint64) uint64 {
    return b / 1024 / 1024
	}
  runtime.ReadMemStats(&m)
  log.Printf("MEM [%s]: Alloc= %v, TotalAlloc= %v, Sys= %v, NumGC = %v",
    s, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC )
}