package config

import (
	"os"

	"github.com/hjson/hjson-go/v4"
	"github.com/rs/zerolog/log"
)

type OrderList []string

// Загрузка одного списка OrderList, по имени файла
func Order_Load(filename string) OrderList {
	// Пустая структура
	target_list := OrderList{}

	// Загрузка данных
	jsonData, err := os.ReadFile(filename)
	if err == nil {
		err = hjson.Unmarshal(jsonData, &target_list)
	}
	if err != nil {
		log.Err(err).Msg(filename + ": load error!")
	}

	return target_list
}


// Сохранение одного списка OrderList, по имени файла
func Order_Save(filename string, target_list OrderList) {
	// Пишем в файл
	jsonData, err := hjson.Marshal(&target_list)
	if err == nil {
		err = os.WriteFile(filename, jsonData, 0664)
		if err == nil {
			log.Info().Msg(filename + ": saved")
			return
		}
	}
	log.Err(err).Msg(filename + ": cannot save")
}
