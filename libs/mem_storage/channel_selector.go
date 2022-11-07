package mem_storage

import (
	"bytes"

	"github.com/rs/zerolog/log"
  
	"ottplay-foss-server-epg/libs/helpers"
)


func ParseEpg_LookupChannel(in []byte, prov_list []*ProviderEpgData, prov_user_len uint8) ([]byte, *EpgChannelData, *ProviderEpgData) {
  var (
    ok bool
    err error
    ch_key []byte
    _ch_key, ch_tid, ch_tname, ch_name uint32
    ch   *EpgChannelData
    prov *ProviderEpgData
    search_step int = 0
  )

  // Ordered будет лежать в другом блоке памяти, удобно :)
  // UPD: Не актуально, но пусть будет
  //      тк, передаем prov_user_len, показывающий "границу" предпочтительных источников
  //std_prov := (&prov_list[0] == &PO.Epg[0])


  cdata := bytes.Split(in, c_sep_vars)  // "-"

  if len(cdata) != 4 { goto err_bye }
    _ch_key, _, err = helpers.ParseUint32Buf(cdata[0])
      if (err != nil) || (_ch_key == 0) { goto err_bye }
      ch_key = cdata[0]
    ch_tid, _, err = helpers.ParseUint32Buf(cdata[1])
      if err != nil { goto err_bye }
    ch_tname, _, err = helpers.ParseUint32Buf(cdata[2])
      if err != nil { goto err_bye }
    ch_name, _, err = helpers.ParseUint32Buf(cdata[3])
      if err != nil { goto err_bye }

  if (ch_tid != 0) && (prov_user_len > 0) {
    // Если список "кастомный", и есть tvg-id
    for i := 0; i < int(prov_user_len); i++ {
      prov = prov_list[i]
      ch, ok = prov.ById[ch_tid]
      if ok { search_step = 1; goto exit_ok }
    }
  }
  if (ch_tname != 0) {
    // Если есть tvg-name
    for i := 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_tname]
      if ok && !ch.ExpiredEpg { search_step = 2; goto exit_ok }
    }
  }
  if (ch_name != 0) {
    // Если есть имя канала
    for i := 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_name]
      if ok && !ch.ExpiredEpg { search_step = 3; goto exit_ok }
    }
  }
  // if (ch_tid != 0) {
  // // Если мы тут, то ничего не осталось, как сверить tvg-id (пока пусть будет TODO)
  // }
  // log.Printf("EPG.NOT_FOUND: %s", ch_key);
  _ = search_step
exit_ok:
  // log.Printf("EPG.FOUND: %s -- %s sate: %d", ch_key, ch.IdHash, search_state);
  return ch_key, ch, prov
err_bye:
  log.Err(err).Msgf("match-channels.epg: cannot parse channel data %s", helpers.B2s(in))
  return nil, nil, nil
}


func ParseIco_LookupChannel(in []byte, prov_list []*ProviderIcoData, prov_user_len uint8) ([]byte, *string, *ProviderIcoData) {
  var (
    ok bool
    err error
    ch_key []byte
    _ch_key, ch_tid, ch_tname, ch_name uint32
    ch   *string
    prov *ProviderIcoData
    search_state int = 0
  )

  // Ordered будет лежать в другом блоке памяти, удобно :)
  // UPD: Не актуально, но пусть будет
  //      тк, передаем prov_user_len, показывающий "границу" предпочтительных источников
  //std_prov := (&prov_list[0] == &PO.Epg[0])

  cdata := bytes.Split(in, c_sep_vars)  // "-"

  if len(cdata) != 4 { goto err_bye }
    _ch_key, _, err = helpers.ParseUint32Buf(cdata[0])
      if (err != nil) || (_ch_key == 0) { goto err_bye }
      ch_key = cdata[0]
    ch_tid, _, err = helpers.ParseUint32Buf(cdata[1])
      if err != nil { goto err_bye }
    ch_tname, _, err = helpers.ParseUint32Buf(cdata[2])
      if err != nil { goto err_bye }
    ch_name, _, err = helpers.ParseUint32Buf(cdata[3])
      if err != nil { goto err_bye }

  if (ch_tid != 0) && (prov_user_len > 0) {
    // Если список "кастомный", и есть tvg-id
    for i := 0; i < int(prov_user_len); i++ {
      prov = prov_list[i]
      ch, ok = prov.ById[ch_tid]
      if ok { search_state = 1; goto exit_ok }
    }
  }
  if (ch_tname != 0) {
    // Если есть tvg-name
    for i := 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_tname]
      if ok { search_state = 2; goto exit_ok }
    }
  }
  if (ch_name != 0) {
    // Если есть имя канала
    for i := 0; i < len(prov_list); i++ {
      prov = prov_list[i]
      ch, ok = prov.ByName[ch_name]
      if ok { search_state = 3; goto exit_ok }
    }
  }
  // if (ch_tid != 0) {
  // // Если мы тут, то ничего не осталось, как сверить tvg-id (пока пусть будет TODO)
  // }
  // log.Printf("ICO.NOT_FOUND: %s", ch_key);
  _ = search_state
  return ch_key, nil, nil
exit_ok:
  // log.Printf("ICO.FOUND: %s -- sate: %d", ch_key, search_state);
  return ch_key, ch, prov
err_bye:
  log.Err(err).Msgf("match-channels.ico: cannot parse channel data %s", helpers.B2s(in))
  return nil, nil, nil
}