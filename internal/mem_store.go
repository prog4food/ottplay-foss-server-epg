package internal

import (
	"sync"

	"ottplay-foss-server-epg/pkg/mem_storage"
)

var (
	ChDbLock sync.RWMutex               // Блокировка на время обновления провайдера
	P        mem_storage.ProvStore      // Глобальное хранилище провайдеров
	PO       mem_storage.OrderedStorage // Отсортированное хранилище провайдеров
)
