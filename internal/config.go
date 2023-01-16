package internal

import "ottplay-foss-server-epg/pkg/config"

var (
	Config   *config.ConfigFile
	OrderEpg config.OrderList
	OrderIco config.OrderList
)

// Загрузка списков OrderList
func LoadOrderLists() {
	OrderEpg = config.Order_Load(config.FileOrderEpg)
	OrderIco = config.Order_Load(config.FileOrderIco)
}
