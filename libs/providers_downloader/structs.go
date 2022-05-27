package providers_downloader

type ProvMeta struct {
  Id         string  `json:"id"`
  Urls     []uint32  `json:"url-hashes"`
  LastUpd    uint64  `json:"last-upd"`
  LastEpg    uint64  `json:"last-epg"`
}

//easyjson:json
type ProviderData struct {
 Meta ProvMeta             `json:"meta"`
 Data map[uint32][]string  `json:"data"`
}