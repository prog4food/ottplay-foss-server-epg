package main

import (
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/internal"
	"ottplay-foss-server-epg/pkg/config"
)

// Устанавливаются при сборке
var depl_ver = "[devel]"

func main() {
	var sOut io.Writer

	// Фикс цветной консоли для Windows
	if runtime.GOOS == "windows" {
		sOut = colorable.NewColorableStdout()
	} else {
		sOut = os.Stdout
	}
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: sOut, TimeFormat: "2006-01-02T15:04:05"})
	//f_client.InitClient()
	log.Info().Msg("OTT-play FOSS server (EPG) " + depl_ver)
	log.Info().Msg("  git@prog4food (c) 2o22")


  internal.Config = config.Load()
  if internal.Config == nil {
    log.Fatal().Msg("config: App could not start with bad config!")
    return
  }

  // Async: Запуск планировщика
	go internal.StartSched()
	// Async: Загрузка конфигурации провайдеров
	go internal.ReLoadConfig()

	log.Info().Msgf("server: starting... [%s]", internal.Config.Bind)

	http.HandleFunc("/html/", internal.ServeStatic)
	http.HandleFunc("/m3u/",  internal.RouterPublic)
	http.HandleFunc("/api/",  internal.PrivateApi)

	err := http.ListenAndServe(internal.Config.Bind, nil)
	if err != nil {
		log.Fatal().Err(err)
	}
}
