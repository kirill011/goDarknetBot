package bot

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"goDark/internal/config"

	"goDark/internal/detector"

	tele "gopkg.in/telebot.v3"
)

// Структура для бота
type Telebot struct {
	n   *detector.Network //переменная в которой хранится экземпляр нейросети
	cfg *config.Config    //переменная в которой лежит структура конфига
}

// функция для создания структуры бота
func NewBot(cfg *config.Config) (*Telebot, error) {
	network, err := detector.Init(cfg) //инициализируем нейросеть
	if err != nil {
		return nil, err //если возникла ошибка, то возвращаем её
	}
	return &Telebot{n: network, cfg: cfg}, nil //возвращаем структуру
}

// инициализация бота
func (t Telebot) Init() {

	//настройки бота
	pref := tele.Settings{
		Token:  t.cfg.ApiKey,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	//создание бота
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		// КОнструктор для меню
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		// создаём переменные кнопок
		btnHelp   = menu.Text("ℹ Помощь")
		btnDetect = menu.Text("👁 Детектировать")
	)

	//отрисовываем кнопки
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnDetect),
	)

	//обработчик кнопки help
	b.Handle(&btnHelp, func(ctx tele.Context) error {
		ctx.Send("Просто отправьте мне фотографию(как фотографию, а не как файл), и я пришлю вам то, что у меня получилось найти. Детектирование происходит с помощью фреймворка darknet(модель yolov3-tiny)")
		for _, val := range t.n.N.ClassNames {
			ctx.Send(val)
		}
		return ctx.Send("Я умею находить только эти классы")
	})

	//обработчик команды start
	b.Handle("/start", func(ctx tele.Context) error {
		return ctx.Send("нажмите на кнопку помощь для получения подсказки", menu)
	})
	//обработчик кнопки detect
	b.Handle(&btnDetect, func(ctx tele.Context) error {
		//добавляем обработчик для детектирования фотографии
		b.Handle(tele.OnPhoto, t.photoHandler)
		return ctx.Send("Отправьте фотографию для детектирования")
	})

	//запуск бота
	b.Start()
}

// обработчик фотографии
func (t *Telebot) photoHandler(ctx tele.Context) error {

	t.cfg.InfoLog.Println("Начинаю детектирование")
	ctx.Send("Начинаю детектирование")
	//получаем фото
	photo := ctx.Message().Photo
	//качаем фото и сохраняем его в папку cache
	photoPath := fmt.Sprintf("internal/cache/%s.jpg", photo.FileID)
	if err := ctx.Bot().Download(&photo.File, photoPath); err != nil {
		t.cfg.ErrorLog.Println(err)
		ctx.Send("Ошибка " + err.Error())
	}

	//открываем скачанное фото
	infile, err := os.Open(photoPath)
	if err != nil {
		t.cfg.ErrorLog.Println(err)       //логируем ошибку
		ctx.Send("Ошибка " + err.Error()) //отправляем ошибку
	}
	defer infile.Close() //после завершения функции закрываем файл

	//декодируем jpeg в image.Image
	src, err := jpeg.Decode(infile)
	if err != nil {
		t.cfg.ErrorLog.Println(err)
		ctx.Send("Ошибка " + err.Error())
	}

	//Детектируем
	count, err := t.n.Detect(t.cfg, &src, photo.FileID+"_detected")
	if err != nil {
		t.cfg.ErrorLog.Println(err)
		ctx.Send("Ошибка " + err.Error())
	}
	//загружаем файл с результатами детектирования из папки cache
	detectedPhoto := tele.FromDisk(fmt.Sprintf("internal/cache/%s_detected.jpg", photo.FileID))
	detected := &tele.Photo{File: detectedPhoto}

	//отправляем сообщение о количестве найденных объектов
	ctx.Send(fmt.Sprintf("Раснознано %d объектов", count))

	t.cfg.InfoLog.Println("Закончил детектирование")
	//отправляем результат детектирования
	return ctx.Send(detected)
}
