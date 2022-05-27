package config_epg

import (
	"encoding/json"
	"os"

	xxhash32 "github.com/OneOfOne/xxhash"
	"github.com/rs/zerolog/log"
)

const config_file = "config.json"
var ConfigData = ConfigFile{}

func Load() {
  jsonData, err := os.ReadFile(config_file); if err != nil {
    log.Fatal().Msg("Cannot load " + config_file)
    return
  }
  
  if err := json.Unmarshal(jsonData, &ConfigData); err != nil {
    log.Err(err).Send()
  }

  // Параметры по-умолчанию
  for i := 0; i < len(ConfigData.Providers); i++ {
    p := &ConfigData.Providers[i] 
    // xxhash32 для провайдера
    p.IdHash = xxhash32.ChecksumString32(p.IdName)
    // Обновление списка каналов = 6 часов
    if p.UpdMins == 0 { p.UpdMins = 360 }
    // Порядок = 100 + порядок элемента
    _default_order := 100 + uint8(i)
    if _default_order > 255 {  _default_order = 255 }
    if p.OrderEpg == 0 { p.OrderEpg = _default_order }
    if p.OrderIco == 0 { p.OrderIco = _default_order }
    
  }
}