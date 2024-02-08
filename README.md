# goDarknetBot

## Описание 
Telegram-бот для детектирования объектов на изображении. Детектирование производится с использованием фреймворка Darknet и модели yolov3-tiny.
Поддерживаются следующие категории объектов:
+ bicycle
+ car
+ motorcycle
+ airplane
+ bus
+ train
+ truck
+ boat

## Запуск 
+ Отредактируйте apiKey и другие параметры в файле ../internal/config/config.yaml
+ Запустите 
  ``` cmd
  go run main.go
  ```
  из директории с проектом
