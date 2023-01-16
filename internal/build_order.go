package internal

import 	(
	"github.com/rs/zerolog/log"
	xxhash32 "github.com/OneOfOne/xxhash"

	"ottplay-foss-server-epg/pkg/mem_storage"
	"ottplay-foss-server-epg/pkg/config"
	"ottplay-foss-server-epg/pkg/helpers"
)

func _sort(listType uint8) {
	var (
		storeId          string
		orderConfig      config.OrderList
		addToGlobalOrder func(prov *mem_storage.ProviderElement)
		// reuse vars
		i                int
		provName         string
		provNameHash     uint32
		ok, orderChanged bool
		provider         *mem_storage.ProviderElement
	)

	if listType == mem_storage.TypeEpg { // Набор переменных для работы с Epg
		storeId = "epg"
		orderConfig = OrderEpg
		PO.Epg = make([]*mem_storage.ProviderEpgData, 0, len(P))
		addToGlobalOrder = func(prov *mem_storage.ProviderElement) {
			if prov.Epg != nil { // nil - значит пропускаем
				if len(prov.Epg.ByName) > 0 {
					PO.Epg = append(PO.Epg, prov.Epg)
				} else {
					log.Error().Msgf("store.%s: %s - empty", storeId, prov.Id)
				}
			}
		}
	} else if listType == mem_storage.TypeIco { // Набор переменных для работы с Ico
		storeId = "ico"
		orderConfig = OrderIco
		PO.Ico = make([]*mem_storage.ProviderIcoData, 0, len(P))
		addToGlobalOrder = func(prov *mem_storage.ProviderElement) {
			if prov.Ico != nil { // nil - значит пропускаем
				if len(prov.Ico.ByName) > 0 {
					PO.Ico = append(PO.Ico, prov.Ico)
				} else {
					log.Error().Msgf("store.%s: %s - empty", storeId, prov.Id)
				}
			}
		}
	}

	// Для сортировки используем временный флаг FlagReserved1, сбрасываем его
	for _, provider = range P {
		provider.Flags = helpers.Bit_clear(provider.Flags, mem_storage.FlagReserved1)
	}

	var _orderLen = len(orderConfig)
	var newOrderConfig = make(config.OrderList, 0, _orderLen+5) // С небольшим запасом на новые элементы
	// Сортировка по orderConfig
	for i = 0; i < _orderLen; i++ {
		provName = orderConfig[i]
		// Все, что начинается с ";" считается комментарием
		if provName[0] != ';' {
			provNameHash = xxhash32.ChecksumString32(provName)
			provider, ok = P[provNameHash] // ищем провайдера в глобальном списке
			if ok {
				if !helpers.Bit_has(provider.Flags, mem_storage.FlagReserved1) {
					// еще НЕ обработан
					addToGlobalOrder(provider)
					provider.Flags = helpers.Bit_set(provider.Flags, mem_storage.FlagReserved1)
				} else {
					// УЖЕ был обработан (дубликат в orderConfig)
					log.Warn().Msgf("sort.%s: %s - duplicate in order", storeId, provName)
					orderChanged = true
					continue
				}
			} else {
				// провайдера НЕТ в глобальном списке
				log.Warn().Msgf("sort.%s: %s - unconfigured", storeId, provName)
			}
		}
		// добавляем в Order на сохранение
		newOrderConfig = append(newOrderConfig, provName)
	}

	// Все новые провайдеры добавляем в конец (в порядке из конфига)
	for i = 0; i < len(Config.Providers); i++ {
		provName = Config.Providers[i].IdName
		provNameHash = xxhash32.ChecksumString32(provName)
		provider, ok = P[provNameHash]
		if !ok {
			log.Error().Msgf("sort.%s: %s - has no mem_storage.ProviderElement!", storeId, provName)
			continue
		}
		if helpers.Bit_has(provider.Flags, mem_storage.FlagReserved1) {
			// Пропускаем ordered и не загруженные источники
			continue
		}
		addToGlobalOrder(provider)
		newOrderConfig = append(newOrderConfig, provider.Id)
		orderChanged = true
	}

	// Пересохраняем список, если требуется
	if orderChanged {
		if listType == mem_storage.TypeEpg {
			OrderEpg = newOrderConfig
			config.Order_Save(config.FileOrderEpg, newOrderConfig)
		} else if listType == mem_storage.TypeIco {
			OrderIco = newOrderConfig
			config.Order_Save(config.FileOrderIco, newOrderConfig)
		}
	}
}

// Пересоздает хранилища Order, с учетом OrderList
func ReOrder() {
	_sort(mem_storage.TypeEpg)
	_sort(mem_storage.TypeIco)
}
