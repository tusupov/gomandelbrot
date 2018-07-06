package mandelbrot

import (
	"crypto/md5"
	"encoding/hex"
"image"
"image/color"
"image/draw"
"math"
	"strconv"
	"sync/atomic"







"github.com/tusupov/gomandelbrot/cache"
)

var inQueueCnt = make([]int32, len(sizesList) + 1)

type mandelbrot struct {
	id					string
	moveX, moveY 		float64
	zoom				uint64
	width, height		int
	priority			int
	param				param
}

type param struct {
	minReal, maxReal 	float64
	minImag, maxImag 	float64
}

// move mandelbrot to x and y pixel
func (m mandelbrot) Move(x, y float64) mandelbrot {

	m.moveX += x
	m.moveY += y * -1

	m.updateId()

	return m

}

// zoom mandelbrot x times
func (m mandelbrot) Zoom(x uint64) mandelbrot {

	if x == 0 { return m }

	m.param = param{
		minReal: m.param.minReal / float64(x),
		maxReal: m.param.maxReal / float64(x),
		minImag: m.param.minImag / float64(x),
		maxImag: m.param.maxImag / float64(x),
	}

	m.zoom *= x

	m.moveX *= float64(x)
	m.moveY *= float64(x)

	m.updateId()

	return m

}

// Добавляет в очередь
func (m *mandelbrot) pushQueue() {

	// Проверка приоритета на валидност
	// если не валидно устанавливает самый низкий приоритет
	if m.priority > len(sizesList) || m.priority < 1 {
		m.priority = len(sizesList)
	}

	// Добавляет в очередь
	for i := m.priority + 1; i <= len(sizesList); i++ {
		atomic.AddInt32(&inQueueCnt[i], 1)
	}

	// Ждет своей очереди
	workers <- struct{}{}

	// Убирает из очереди
	for i := m.priority + 1; i <= len(sizesList); i++ {
		atomic.AddInt32(&inQueueCnt[i], -1)
	}

}

// Убирает из очереди
func (m mandelbrot) pullQueue() {
	<- workers
}

// Проверка есть ли более приоритетные процессы
// которые попали в очередь
func (m mandelbrot) hasLowerInQueue() bool {

	if inQueueCnt[m.priority] > 0 {
		return true
	}

	return false
}

func (m mandelbrot) Draw() (buf []byte) {

	// Загружает из очереди
	buf, err := cache.Load(m.width, m.height, m.id)
	if err == nil {
		return
	}

	// Добавляет и убирает из очереди
	m.pushQueue()
	defer m.pullQueue()

	// Создает новую картинку с указанными размерами
	img := image.NewRGBA64(image.Rect(0, 0, m.width, m.height))

	// Заполняет картинку черным цветом
	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{color.Black},
		img.Bounds().Min,
		draw.Over,
	)

	// Количество итерации
	iteration := m.calcIteration()

	// Максимальны значение цветом
	maxuint16 := ^uint16(0)
	colorK := maxuint16 / uint16(iteration)

	for x := 0; x <= m.width; x++ {

		for y := 0; y <= m.height; y++ {

			c := m.complex(x, y)
			z := complex(0, 0)

			for i := 0; i < iteration; i++ {

				z = z*z + c

				if real(z)*real(z)+imag(z)*imag(z) > 4 {

					// Определяет нужный оттенок
					clr := uint16(i) * colorK

					// Красит точку определенным цветом
					// черно белым цветом
					img.Set(
						x,
						y,
						color.RGBA64{
							R: clr,
							G: clr,
							B: clr,
							A: maxuint16,
						},
					)

					break

				}

			}

		}

		// Проверка есть ли более приоритетные процессы
		// которые попали в очередь
		if m.hasLowerInQueue() {
			m.pullQueue()
			m.pushQueue()
		}

	}

	// Сохраняет в кэш и возвращаем байт массив картинки
	buf, err = cache.Save(m.width, m.height, m.id, img)

	return

}

// Рассчитывает complex `c` с параметрами `x` and `y`
func (m mandelbrot) complex(x, y int) complex128 {

	return complex(
		(m.param.maxReal - m.param.minReal) * (float64(x) + m.moveX) / float64(m.width) + m.param.minReal,
		(m.param.maxImag - m.param.minImag) * (float64(y) + m.moveY) / float64(m.height) + m.param.minImag,
	)

}

// Рассчитывает количество итерации
// в зависимости от увеличение
func (m mandelbrot) calcIteration() int {

	f := math.Sqrt(
		0.001 +
		2.0 * math.Min(
			math.Abs(m.param.minReal - m.param.maxReal),
			math.Abs(m.param.minImag - m.param.maxImag),
		),
	)

	return int(208.0 / f)

}

// Устанавливает ид для этого экземпляра
func (m *mandelbrot) updateId() {

	str := strconv.FormatInt(int64(m.zoom), 10) + "_" +
			strconv.FormatFloat(m.moveX, 'f', -1, 64) + "_" +
			strconv.FormatFloat(m.moveY, 'f', -1, 64)

	hasher := md5.New()
	hasher.Write([]byte(str))

	m.id = hex.EncodeToString(hasher.Sum(nil))

}

// Создаем новый экземпляр
func New(rec string) mandelbrot {

	// Получат значение размера
	width, height, priority := GetSize(rec)

	m := mandelbrot{
		width:		width,
		height:		height,
		priority:	priority,
		zoom:		1,
		param:		param{
			minReal: 	-2.0,
			maxReal: 	1.0,
			minImag:	-1.5,
			maxImag:	1.5,
		},
	}

	// Устанавливает ид для этого экземпляра
	m.updateId()

	return m

}
