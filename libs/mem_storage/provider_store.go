package mem_storage

import (
	"sort"
	"sync"
)

// Хранилище провайдера Epg
type EpgChannelData struct {
  IdHash    []byte
  // // TODO: пока нет смысла сохранять имена и Id в виде строк
  // Id         *string
  // Names    []*string
  LastEpg     int64
  ExpiredEpg  bool
}
type HashEpgStruct map[uint32]*EpgChannelData
type ProviderEpgData struct {
  ById        HashEpgStruct
  ByName      HashEpgStruct
  Order       uint16
  UrlBase    *string
  UrlBaseS   *string
  Outdated    uint32
  Prov       *ProviderElement
}
func (m ProviderEpgData) GetOrder() uint16 { return m.Order }
type StoreTypes interface {
  *ProviderEpgData | *ProviderIcoData
  GetOrder() uint16
}

// Хранилище провайдера Ico
type HashIcoStruct map[uint32]*string
type ProviderIcoData struct {
  ById    HashIcoStruct
  ByName  HashIcoStruct
  Order   uint16
  Prov   *ProviderElement
}
func (m ProviderIcoData) GetOrder() uint16 { return m.Order }

// Единица хранилища провайдера
type ProviderElement struct {
  Id         *string
  IdHashes  []uint32
  Epg        *ProviderEpgData
  Ico        *ProviderIcoData
}

type OrderedStorage struct {
  Epg  []*ProviderEpgData
  Ico  []*ProviderIcoData
}
var (
  Lock sync.RWMutex
  P = make(map[uint32]*ProviderElement) // Общее map провайдеров
  PO OrderedStorage                // Общее сортированное хранилище провайдеров
)


// Сортировка хранилищ провайдеров Epg/Ico
func (t *OrderedStorage) Sort() {
  e_slice := make([]*ProviderEpgData, 0, len(P))
  for _, v := range P { e_slice = append(e_slice, v.Epg) }
  t.Epg = StoreSlice_Sort(e_slice)

  i_slice := make([]*ProviderIcoData, 0, len(P))
  for _, v := range P { i_slice = append(i_slice, v.Ico) }
  t.Ico = StoreSlice_Sort(i_slice)
}


// Дженерик сортировки провайдеров по их свойству Order
func StoreSlice_Sort[T StoreTypes](store []T, ) []T {
  sort.Slice(store, func(i, j int) bool {
    return store[i].GetOrder() < store[j].GetOrder()
  })
  return store
}
