package detector

import (
	"fmt"
	"goDark/internal/config"
	"image"
	"image/color"
	"log"
	"strconv"

	"github.com/LdDl/go-darknet"
	"github.com/fogleman/gg"
)

// структура пакета
type Network struct {
	N *darknet.YOLONetwork
}

// инициализация нейросети
func Init(cfg *config.Config) (*Network, error) {
	n := darknet.YOLONetwork{
		GPUDeviceIndex:           0,               //индекс видеокарты
		NetworkConfigurationFile: cfg.ConfigPath,  //путь к конфигурации нейросети
		WeightsFile:              cfg.WeightsPath, // путь к файлу весов
		Threshold:                cfg.Threshold,   // порог ниже которого продетектированный объект не попадает в результаты
	}
	//если во время инициализации нейросети (инициализация из внешнего пакета) возникла ошибка
	if err := n.Init(); err != nil {
		return nil, err
	}

	return &Network{N: &n}, nil
}

// функция детектирования
func (n *Network) Detect(cfg *config.Config, img *image.Image, fileName string) (int, error) {

	//преобразуем фото в формат darknet
	imgDarknet, err := darknet.Image2Float32(*img)
	if err != nil {
		return 0, err
	}

	//детектируем
	dr, err := n.N.Detect(imgDarknet)
	if err != nil {
		return 0, err
	}
	imgDarknet.Close()

	log.Println("Network-only time taken:", dr.NetworkOnlyTimeTaken)
	log.Println("Overall time taken:", dr.OverallTimeTaken, len(dr.Detections))

	//создаём контекст изображения
	ctx := gg.NewContextForImage(*img)

	//проходимся по массиву продетектированных объектов
	for _, d := range dr.Detections {
		for i := range d.ClassIDs {
			//выводим в консоль результаты детектирования
			bBox := d.BoundingBox
			fmt.Printf("%s (%d): %.4f%% | start point: (%d,%d) | end point: (%d, %d)\n",
				d.ClassNames[i], d.ClassIDs[i],
				d.Probabilities[i],
				bBox.StartPoint.X, bBox.StartPoint.Y,
				bBox.EndPoint.X, bBox.EndPoint.Y,
			)

			//создаём подпись над рамкой объекта
			title := d.ClassNames[i] + " " + strconv.FormatFloat(float64(d.Probabilities[i]), 'f', 4, 32) + "%"

			// получаем минимумы и максимумы рамки продетектированного объекта
			minX, minY := float64(bBox.StartPoint.X), float64(bBox.StartPoint.Y)
			maxX, maxY := float64(bBox.EndPoint.X), float64(bBox.EndPoint.Y)

			//рисуем прямоугольник по минимумам и максимумам
			ctx.DrawRectangle(minX, minY, maxX-minX, maxY-minY)
			//выставляем ширину рамки
			ctx.SetLineWidth(2.0)
			//загружаем шрифт для title
			ctx.LoadFontFace(cfg.FontPath, 24)

			//в зависимости от того, какой объект мы обнаружили, выставляем цвета рамки и title
			switch d.ClassNames[i] {
			case "bicycle":
				ctx.SetColor(color.NRGBA{255, 0, 0, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "car":
				ctx.SetColor(color.NRGBA{255, 255, 0, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "motorcycle":
				ctx.SetColor(color.NRGBA{255, 255, 255, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "airplane":
				ctx.SetColor(color.NRGBA{0, 255, 0, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "bus":
				ctx.SetColor(color.NRGBA{0, 255, 255, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "train":
				ctx.SetColor(color.NRGBA{0, 0, 255, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "truck":
				ctx.SetColor(color.NRGBA{255, 0, 255, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			case "boat":
				ctx.SetColor(color.NRGBA{0, 0, 0, 100})
				ctx.DrawString(title, minX, minY)
				ctx.Stroke()
				break
			}
		}
	}

	//сохраняем в cache
	ctx.SavePNG(fmt.Sprintf("internal/cache/%s.jpg", fileName))

	return len(dr.Detections) - 1, nil //возвращаем количество продетектированных объектов
}
