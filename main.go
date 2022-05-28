package main

import (
	_ "net/http/pprof"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"ottplay-foss-server-epg/libs/config_epg"
	"ottplay-foss-server-epg/libs/providers_downloader"
	"ottplay-foss-server-epg/libs/web_logic"
)

// Устанавливаются при сборке
var depl_ver string


var (
  c_byte_wild = []byte{'*'}
)

func fastHTTPError(ctx *fasthttp.RequestCtx, ecode int) {
  ctx.SetStatusCode(ecode)
}

func AllowCors(ctx *fasthttp.RequestCtx) {
  ctx.Response.Header.AddBytesV(fasthttp.HeaderAccessControlAllowOrigin, c_byte_wild)
}

func main() {
  //zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
  log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05", NoColor: true})
  //f_client.InitClient()
  log.Info().Msg("OTT-play FOSS server (EPG) " + depl_ver)
  log.Info().Msg("  git@prog4food (c) 2o22")

  // Загрузка конфигурации
  config_epg.Load()
  go providers_downloader.StartJob()

  // the corresponding fasthttp code
  m := func(ctx *fasthttp.RequestCtx) {
    switch string(ctx.Path()) {
    case "/m3u/match-logos":
      web_logic.IcoMatch(ctx)
      AllowCors(ctx)
    case "/m3u/match-channels":
      web_logic.EpgMatch(ctx)
      AllowCors(ctx)
    case "/api/update-provider":
      // TODO: update hook
    case "/m3u/gelist.php":
    case "/m3u/geicons.php":
      //TODO: backward compability
      //TODO: backward compability
    default:
      ctx.Error("not found", fasthttp.StatusNotFound)
    }
  }

  fasthttp.ListenAndServe(config_epg.ConfigData.Bind, m)
}
