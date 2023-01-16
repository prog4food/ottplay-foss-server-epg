package request_parser

import (
	"bytes"
	"net/http"

	"ottplay-foss-server-epg/pkg/mem_storage"
)

var (
	c_sep_block  = []byte("\n\t\n") // \n\t\n
	c_sep_line   = []byte{0x0A}     // \n
	c_sep_array  = []byte{'~'}      // ~
	c_sep_vars   = []byte{'-'}      // -
	c_empty_meta = []byte("{}")     // {}
)

func ParseEpg_ReqChannels(in []byte, prov_list []*mem_storage.ProviderEpgData, w http.ResponseWriter, prov_user_len uint8) {
	if len(in) == 0 {
		return
	}
	clist := bytes.Split(in, c_sep_line) // "\n"

	var providers []*mem_storage.ProviderEpgData
	var cdata [][]byte
	var (
		// Блок переменных для используемых провайдеров
		used_provs     = make([]*mem_storage.ProviderEpgData, 0, len(prov_list))
		used_provs_len int
		// Блок переменных для результатов поиска канала
		ch_key  []byte
		ch_data *mem_storage.EpgChannelData
		ch_prov *mem_storage.ProviderEpgData
	)

	// Пишем json заголовок в файл
	//out_buf  := &ctx.Response
	mini_buf := make([]byte, 0, 256)

	// Блок: МЕТА
	w.Write(c_empty_meta) // TODO

	// Блок: КАНАЛЫ
	w.Write(c_sep_block)
	// Проход по каналам
	for i := 0; i < len(clist); i++ {
		if len(clist[i]) == 0 {
			continue
		}

		cdata = bytes.Split(clist[i], c_sep_array) // "~" - делим данные канала и его epg
		if len(cdata) == 2 {
			// Приоритизируем пользовательские источники
			providers, prov_user_len = ReadChannelProviders(cdata[1], prov_list, Epg_GetProvByHash)
		} else {
			// Используем общий список
			providers = prov_list
		}

		ch_key, ch_data, ch_prov = ParseAndLookup_Epg(cdata[0], providers, prov_user_len)
		if ch_data != nil {
			// Пишем данные о канале
			mini_buf = mini_buf[:0] // Сброс буфера
			used_provs = Slice_AppendUniq(used_provs, ch_prov)
			if i > 0 {
				mini_buf = append(mini_buf, c_sep_line...)
			}
			mini_buf = append(mini_buf, ch_key...)
			mini_buf = append(mini_buf, c_sep_array...)
			mini_buf = append(mini_buf, ch_prov.Prov.Id...)
			mini_buf = append(mini_buf, c_sep_array...)
			mini_buf = append(mini_buf, ch_data.IdHash...)
			w.Write(mini_buf)
		}
	}
	w.Write(c_sep_block)
	used_provs_len = len(used_provs)

	// Блок: ПРОВАЙДЕРЫ
	if used_provs_len > 0 {
		// Пишем данные об используемых провайдерах
		for i := 0; i < used_provs_len; i++ {
			mini_buf = mini_buf[:0] // Сброс буфера
			if i > 0 {
				mini_buf = append(mini_buf, c_sep_line...)
			}
			mini_buf = append(mini_buf, used_provs[i].Prov.Id...)
			mini_buf = append(mini_buf, c_sep_array...)
			mini_buf = append(mini_buf, used_provs[i].ProvUrl...)
			w.Write(mini_buf)
		}
	}
}

func ParseIco_ReqChannels(in []byte, prov_list []*mem_storage.ProviderIcoData, w http.ResponseWriter, prov_user_len uint8) {
	if len(in) == 0 {
		return
	}
	clist := bytes.Split(in, c_sep_line) // "\n"

	var providers []*mem_storage.ProviderIcoData
	var cdata [][]byte
	var (
		// Блок переменных для результатов поиска канала
		ch_key  []byte
		ch_data []byte
	)
	// Пишем json заголовок в файл
	mini_buf := make([]byte, 0, 256)

	// Блок: МЕТА
	w.Write(c_empty_meta)

	// Блок: КАНАЛЫ
	w.Write(c_sep_block)
	// Проход по каналам
	for i := 0; i < len(clist); i++ {
		if len(clist[i]) == 0 {
			continue
		}

		cdata = bytes.Split(clist[i], c_sep_array) // "~" - делим данные канала и его epg
		if len(cdata) == 2 {
			// Приоритизируем пользовательские источники
			providers, prov_user_len = ReadChannelProviders(cdata[1], prov_list, Ico_GetProvByHash)
		} else {
			// Используем общий список
			providers = prov_list
		}

		ch_key, ch_data, _ = ParseAndLookup_Ico(cdata[0], providers, prov_user_len)
		if ch_data != nil {
			// Пишем данные о канале
			mini_buf = mini_buf[:0] // Сброс буфера
			if i > 0 {
				mini_buf = append(mini_buf, c_sep_line...)
			}
			mini_buf = append(mini_buf, ch_key...)
			mini_buf = append(mini_buf, c_sep_array...)
			mini_buf = append(mini_buf, ch_data...)
			w.Write(mini_buf)
		}
	}
}
