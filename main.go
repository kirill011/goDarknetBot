package main

import (
	"goDark/internal/bot"
	"goDark/internal/config"
	"os"
	"time"
)

func main() {
	//инициализируем конфиг
	cfg := config.Init()

	cfg.InfoLog.Println("Инициализация нейросети")
	bot, err := bot.NewBot(cfg)
	if err != nil {
		cfg.ErrorLog.Fatal(err)
	}

	cfg.InfoLog.Println("Инициализация отчистки кэша")
	go clearCache(cfg)

	cfg.InfoLog.Println("Инициализация бота")
	bot.Init()
}

// функция для отчистки папки cache
func clearCache(cfg *config.Config) {
	for {
		//выбираем все файлы из пвпки
		files, err := os.ReadDir("internal/cache")
		if err != nil {
			cfg.ErrorLog.Println("func app.MonitFiles.ReadDir: ", err)
		}
		//удаляем все выбранные файлы
		for _, file := range files {
			os.Remove("internal/cache/" + file.Name())
		}

		cfg.InfoLog.Println("Кэш отчищен")
		time.Sleep(cfg.CacheCleanup) //ждём заданное время
	}
}
