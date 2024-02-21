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
1. Устапновите [CUDA toolkit](https://developer.nvidia.com/cuda-toolkit) с официального сайта
2. Установите [Darknet](https://github.com/pjreddie/darknet)
3. Установите  [go-darknet](https://github.com/LdDl/go-darknet)
4. Отредактируйте файл go-darknet/darknet.go. Установите пути к дирректориям CUDA и darknet
5. Отредактируйте apiKey и другие параметры в файле ../internal/config/config.yaml
+ Запустите 
  ``` cmd
  go run main.go
  ```
  из директории с проектом
