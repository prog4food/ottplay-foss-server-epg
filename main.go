package main

import (
	"bytes"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"ottplay-foss-server-epg/libs/config_epg"
	"ottplay-foss-server-epg/libs/providers_downloader"
	"ottplay-foss-server-epg/libs/web_logic"
)

// Устанавливаются при сборке
var depl_ver = "[devel]"

var (
  c_byte_wild = []byte{'*'}
  c_byte_static_dir = []byte("/html/")
)

func fastHTTPError(ctx *fasthttp.RequestCtx, ecode int) {
  ctx.SetStatusCode(ecode)
}

func AllowCors(ctx *fasthttp.RequestCtx) {
  ctx.Response.Header.AddBytesV(fasthttp.HeaderAccessControlAllowOrigin, c_byte_wild)
}

func main() {
  var sOut io.Writer

  // Фикс цветной консоли для Windows
  if runtime.GOOS == "windows" {
    sOut = colorable.NewColorableStdout()
  } else {
    sOut = os.Stdout }
  //zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
  log.Logger = log.Output(zerolog.ConsoleWriter{Out: sOut, TimeFormat: "2006-01-02T15:04:05"})
  //f_client.InitClient()
  log.Info().Msg("OTT-play FOSS server (EPG) " + depl_ver)
  log.Info().Msg("  git@prog4food (c) 2o22")
  // Загрузка конфигурации
  config_epg.Load()
  go providers_downloader.StartJob()


  // Локальный хостинг для EPG
  fs := &fasthttp.FS{
    IndexNames:         []string{"index.html"},
    GenerateIndexPages: true,
    Compress:           true,
  }
  html_handler := fs.NewRequestHandler()

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
      // local hosting?
      if bytes.HasPrefix(ctx.Path(), c_byte_static_dir) {
        html_handler(ctx)
        AllowCors(ctx)
      } else {
        ctx.Error("not found", fasthttp.StatusNotFound)
      }
    }
  }

  fasthttp.ListenAndServe(config_epg.ConfigData.Bind, m)
}
