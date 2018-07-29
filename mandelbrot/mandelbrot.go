package mandelbrot

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"math"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

// Workers
var workers = make(chan struct{}, 1)

func SetWorkers(cnt uint) error {

	if cnt == 0 {
		return errors.New("Workers count must not be zero")
	}

	workers = make(chan struct{}, cnt)
	return nil

}

type mandelbrot struct {
	id            string
	moveX, moveY  float64
	zoom          uint64
	width, height int
	priority      int
	param         param
}

type param struct {
	minReal, maxReal float64
	minImag, maxImag float64
}

// Move coordinate
func (m mandelbrot) Move(x, y float64) mandelbrot {

	m.moveX += x
	m.moveY += y * -1
	m.updateId()

	return m

}

// Zoom x times
func (m mandelbrot) Zoom(x uint64) mandelbrot {

	if x == 0 {
		return m
	}

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

// Push in Queue
func (m *mandelbrot) pushQueue() {

	// Check priority
	if m.priority < 0 || len(sizeList) <= m.priority {
		m.priority = len(sizeList)
	}

	// Add to in Queue
	for i := m.priority + 1; i < len(sizeList); i++ {
		atomic.AddInt32(&inQueueSum[i], 1)
	}

	// Wait
	workers <- struct{}{}

	// Remove from in Queue
	for i := m.priority + 1; i < len(sizeList); i++ {
		atomic.AddInt32(&inQueueSum[i], -1)
	}

}

// Pull from Queue
func (m mandelbrot) pullQueue() {
	<-workers
}

func (m mandelbrot) hasLowerInQueue() bool {

	if inQueueSum[m.priority] > 0 {
		return true
	}

	return false
}

func (m mandelbrot) Draw() (buf []byte, err error) {

	m.pushQueue()
	defer m.pullQueue()

	img := image.NewRGBA64(image.Rect(0, 0, m.width, m.height))

	iteration := m.calcIteration()

	// Colors
	maxuint16 := ^uint16(0)
	colorK := maxuint16 / uint16(iteration)

	for x := 0; x <= m.width; x++ {

		for y := 0; y <= m.height; y++ {

			c := m.complex(x, y)
			z := complex(0, 0)

			for i := 0; i < iteration; i++ {

				z = z*z + c

				if real(z)*real(z)+imag(z)*imag(z) > 4 {

					// Calc color
					clr := uint16(i) * colorK

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

		// Check InQueue
		if m.hasLowerInQueue() {
			m.pullQueue()
			m.pushQueue()
		}

	}

	r := new(bytes.Buffer)
	err = png.Encode(r, img)
	if err != nil {
		return
	}

	buf = r.Bytes()

	return

}

// Calculate complex `c` with params `x` and `y`
func (m mandelbrot) complex(x, y int) complex128 {

	return complex(
		(m.param.maxReal-m.param.minReal)*(float64(x)+m.moveX)/float64(m.width)+m.param.minReal,
		(m.param.maxImag-m.param.minImag)*(float64(y)+m.moveY)/float64(m.height)+m.param.minImag,
	)

}

// Calculate iteration
func (m mandelbrot) calcIteration() int {

	f := math.Sqrt(
		0.001 +
			2.0*math.Min(
				math.Abs(m.param.minReal-m.param.maxReal),
				math.Abs(m.param.minImag-m.param.maxImag),
			),
	)

	return int(208.0 / f)

}

// Update Id
func (m *mandelbrot) updateId() {

	str := strconv.FormatInt(int64(m.zoom), 10) + "_" +
		strconv.FormatFloat(m.moveX, 'f', -1, 64) + "_" +
		strconv.FormatFloat(m.moveY, 'f', -1, 64)

	hasher := md5.New()
	hasher.Write([]byte(str))

	m.id = hex.EncodeToString(hasher.Sum(nil))

}

// New mandelbrot
func New(rec string) (m mandelbrot, err error) {

	// Get size
	size, priority, err := GetSize(rec)
	if err != nil {
		return
	}

	m = mandelbrot{
		width:    size,
		height:   size,
		priority: priority,
		zoom:     1,
		param: param{
			minReal: -2.0,
			maxReal: 1.0,
			minImag: -1.5,
			maxImag: 1.5,
		},
	}

	m.updateId()

	return

}
