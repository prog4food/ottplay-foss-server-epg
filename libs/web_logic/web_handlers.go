package web_logic

import (
	"bytes"

	"github.com/valyala/fasthttp"

	ms "ottplay-foss-server-epg/libs/mem_storage"
)

var (
  c_sep_block = []byte("\n\t\n")    // \n\t\n
)

func fastHTTPError(ctx *fasthttp.RequestCtx, ecode int) {
  ctx.SetStatusCode(ecode)
}

func EpgMatch(ctx *fasthttp.RequestCtx) {
  if !ctx.IsPost() {
    fastHTTPError(ctx, 400)
    return
  }
  req_parts := bytes.Split(ctx.Request.Body(), c_sep_block)
  if len(req_parts) != 3 { return }
  // Читаем список провайдеров из заголовка
  // TODO: обработка метаданных из req_parts[1]
  ms.Lock.RLock()
  defer ms.Lock.RUnlock()
  _m3u_providers, prov_user_len := ms.PrioritizeUserProviders(req_parts[1], ms.PO.Epg, ms.Epg_GetProvByHash)
  if len(_m3u_providers) == 0 { fastHTTPError(ctx, 503,); return }
  ms.ParseEpg_ReqChannels(req_parts[2], _m3u_providers, ctx, prov_user_len)
}


func IcoMatch(ctx *fasthttp.RequestCtx) {
  if !ctx.IsPost() {
    fastHTTPError(ctx, 400)
    return
  }
  req_parts := bytes.Split(ctx.Request.Body(), c_sep_block)
  if len(req_parts) != 3 { return }
  // Читаем список провайдеров из заголовка
  // TODO: обработка метаданных из req_parts[1]
  ms.Lock.RLock()
  defer ms.Lock.RUnlock()
  _m3u_providers, prov_user_len := ms.PrioritizeUserProviders(req_parts[1], ms.PO.Ico, ms.Ico_GetProvByHash)
  if len(_m3u_providers) == 0 { fastHTTPError(ctx, 503,); return }
  ms.ParseIco_ReqChannels(req_parts[2], _m3u_providers, ctx, prov_user_len)
}
