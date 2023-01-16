package request_parser

import (
	"bytes"

	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/pkg/helpers"
	"ottplay-foss-server-epg/pkg/mem_storage"
)


func splitChannelData(in []byte) (error, []byte, uint32, uint32, uint32) {
  var err error
  var out [3]uint32
  cdata := bytes.Split(in, c_sep_vars)  // "-"
    if len(cdata) != 4 { goto err_bye }
    out[0], _, err = helpers.ParseUint32Buf(cdata[0])
      if (err != nil) || (out[0] == 0) { goto err_bye }
    out[0], _, err = helpers.ParseUint32Buf(cdata[1])
      if err != nil { goto err_bye }
    out[1], _, err = helpers.ParseUint32Buf(cdata[2])
      if err != nil { goto err_bye }
    out[2], _, err = helpers.ParseUint32Buf(cdata[3])
      if err != nil { goto err_bye }
    return nil, cdata[0], out[0],out[1],out[2]
  err_bye:
    return err, cdata[0], out[0],out[1],out[2]
}

func ParseAndLookup_Epg(in []byte, prov_list []*mem_storage.ProviderEpgData, prov_user_len uint8) ([]byte, *mem_storage.EpgChannelData, *mem_storage.ProviderEpgData) {
  var (
    i   int
    ch *mem_storage.EpgChannelData
    ok  bool
    err error

    ch_key_str []byte
    ch_tid, ch_tname, ch_name uint32
    prov *mem_storage.ProviderEpgData
    search_step int8
  )

  // Ordered будет лежать в другом блоке памяти, удобно :)
  // UPD: Не актуально, но пусть будет
  //      тк, передаем prov_user_len, показывающий "границу" предпочтительных источников
  //std_prov := (&prov_list[0] == &PO.Epg[0])

  err, ch_key_str, ch_tid, ch_tname, ch_name  = splitChannelData(in)
  if err != nil { goto err_bye }

  if (ch_tid != 0) && (prov_user_len > 0) {
    // Если список "кастомный", и есть tvg-id
    for i = 0; i < int(prov_user_len); i++ {
      prov = prov_list[i]
      ch, ok = prov.ById[ch_tid]
      if ok && !ch.ExpiredEpg { search_step = 1; goto exit_ok }
    }
  }
  if (ch_tname != 0) {
    // Если есть tvg-name
    for i = 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_tname]
      if ok && !ch.ExpiredEpg { search_step = 2; goto exit_ok }
    }
  }
  if (ch_name != 0) {
    // Если есть имя канала
    for i = 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_name]
      if ok && !ch.ExpiredEpg { search_step = 3; goto exit_ok }
    }
  }
  // if (ch_tid != 0) {
  // // Если мы тут, то ничего не осталось, как сверить tvg-id (пока пусть будет TODO)
  // }
  // log.Printf("EPG.NOT_FOUND: %s", helpers.B2s(ch_key_str));
  _ = search_step
exit_ok:
  // log.Printf("EPG.FOUND: %s -- %s sate: %d", helpers.B2s(ch_key_str), ch.IdHash, search_step);
  return ch_key_str, ch, prov
err_bye:
  log.Err(err).Msgf("match-channels.epg: cannot parse channel data %s", helpers.B2s(in))
  return nil, nil, nil
}


func ParseAndLookup_Ico(in []byte, prov_list []*mem_storage.ProviderIcoData, prov_user_len uint8) ([]byte, []byte, *mem_storage.ProviderIcoData) {
  var (
    i   int
    ch  []byte
    ok  bool
    err error
    
    ch_key_str []byte
    ch_tid, ch_tname, ch_name uint32
    prov *mem_storage.ProviderIcoData
    search_step int8 
  )

  // Ordered будет лежать в другом блоке памяти, удобно :)
  // UPD: Не актуально, но пусть будет
  //      тк, передаем prov_user_len, показывающий "границу" предпочтительных источников
  //std_prov := (&prov_list[0] == &PO.Epg[0])

  err, ch_key_str, ch_tid, ch_tname, ch_name  = splitChannelData(in)
  if err != nil { goto err_bye }

  if (ch_tid != 0) && (prov_user_len > 0) {
    // Если список "кастомный", и есть tvg-id
    for i = 0; i < int(prov_user_len); i++ {
      prov = prov_list[i]
      ch, ok = prov.ById[ch_tid]
      if ok { search_step = 1; goto exit_ok }
    }
  }
  if (ch_tname != 0) {
    // Если есть tvg-name
    for i = 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_tname]
      if ok { search_step = 2; goto exit_ok }
    }
  }
  if (ch_name != 0) {
    // Если есть имя канала
    for i = 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_name]
      if ok { search_step = 3; goto exit_ok }
    }
  }
  // if (ch_tid != 0) {
  // // Если мы тут, то ничего не осталось, как сверить tvg-id (пока пусть будет TODO)
  // }
  // log.Printf("ICO.NOT_FOUND: %s", ch_key);
  _ = search_step
  return ch_key_str, nil, nil
exit_ok:
  // log.Printf("ICO.FOUND: %s -- sate: %d", ch_key, search_state);
  return ch_key_str, ch, prov
err_bye:
  log.Err(err).Msgf("match-channels.ico: cannot parse channel data %s", helpers.B2s(in))
  return nil, nil, nil
}