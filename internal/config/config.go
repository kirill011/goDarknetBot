package config

import (
	"log"
	"os"
	"time"

	"github.com/kkyr/fig"
)

// структура конфига
type Config struct {
	ApiKey       string        `fig: "apiKey"`      //ключ API
	FontPath     string        `fig: "fontPath"`    //путь к шрифту
	ConfigPath   string        `fig: "configPath"`  //путь к конфигурации нейросети
	WeightsPath  string        `fig: "weightsPath"` //путь к весам
	Threshold    float32       `fig:"threshold`     //порог отсечения
	CacheCleanup time.Duration `fig:"cacheCleanup"` //интервал очистки кэша

	//логеры
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// инициализация конфига
func Init() *Config {
	cfg := Config{}

	//запись логгеров в конфиг
	cfg.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	cfg.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	//загружаем из конфиг файла
	if err := fig.Load(&cfg,
		fig.File("config.yaml"),                                         //имя файла
		fig.Dirs("internal/config", "../internal/config")); err != nil { //путь
		cfg.ErrorLog.Fatalln("func config.Init: ", err) // если ошибка, то логируем
	}

	cfg.InfoLog.Println("Config loaded")
	return &cfg //возвращаем переменную конфига
}
