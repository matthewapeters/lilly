package transform

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// MyKernel -- small matrix for convolving with the image to detect edges
type MyKernel [3][3]float64

// IJ because it is i,j
type IJ struct {
	i int
	j int
}

// ThreeByThree static iteration of nested 3x3
var ThreeByThree = []IJ{
	{0, 0}, {0, 1}, {0, 2},
	{1, 0}, {1, 1}, {1, 2},
	{2, 0}, {2, 1}, {2, 2},
}

// Gx Horizontal Kernel
// var Gx Kernel = [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
var Gx MyKernel = [3][3]float64{
	{-1.0, 0.0, 1.0},
	{-2.0, 0.0, 2.0},
	{-1.0, 0.0, 1.0}}

// Gy Vertical Kernel
// var Gy Kernel = [3][3]int{{1, 2, 1}, {0, 0, 0}, {-1, -2, -1}}
var Gy MyKernel = [3][3]float64{
	{+1.0, +2.0, +1.0},
	{+0.0, +0.0, +0.0},
	{-1.0, -2.0, -1.0}}

type EdgeDetectConfig struct {
	// Reg, Green, Blue factors for adjusting luminance
	RedFactor   float64
	GreenFactor float64
	BlueFactor  float64

	F         float64
	S         float64
	Tx        bool // set background to transparent?
	ShowAngle bool
}

func DefaultEdgeDetectConfig() *EdgeDetectConfig {
	return &EdgeDetectConfig{
		RedFactor:   0.299,
		GreenFactor: 0.587,
		BlueFactor:  0.114,
		F:           127,
		S:           5.0,
		Tx:          true,
	}
}

func computePixelGray(img *image.Gray, x, y int) (float64, float64) {
	window := [3][3]float64{}
	var gradientX, gradientY float64
	gradientX = 0
	gradientY = 0

	for _, ij := range ThreeByThree {
		r := img.GrayAt(x-1+ij.i, y-1+ij.j).Y
		window[ij.i][ij.j] = float64(r)
		gradientX += (window[ij.i][ij.j] * Gx[ij.i][ij.j])
		gradientY += (window[ij.i][ij.j] * Gy[ij.i][ij.j])
	}
	g := (math.Sqrt(float64((gradientX * gradientX) + (gradientY * gradientY))))
	angle := math.Atan2(gradientY, gradientX)

	return g, angle
}

func doARow(y int, x1, x2 int, img *image.Gray, newimg *image.Gray, wg *sync.WaitGroup) {
	defer wg.Done()
	// The convolved value of a single point
	var gradient float64
	for x := x1; x < x2; x++ {
		// compute gradient and ganle
		gradient, _ = computePixelGray(img, x, y)

		var Y uint8 = uint8(gradient)
		newimg.Set(x, y, color.Gray{Y})
	}
}

func ApplySigmoid(img *image.Gray, cfg *EdgeDetectConfig) (newImg *image.RGBA) {
	b := img.Bounds()
	newImg = image.NewRGBA(b)
	var clr color.Color
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lum := uint8(255 / (1 + math.Exp(-1.0*(float64(img.GrayAt(x, y).Y)-cfg.F)/cfg.S)))
			if cfg.Tx && lum < uint8(cfg.F) {
				clr = color.RGBA{0, 0, 0, 0}
			} else {
				clr = color.Gray{lum}
			}
			newImg.Set(x, y, clr)
		}
	}
	return
}

// ImageToGray converts an RGBA image to Gray
func ImageToGray(img *image.RGBA) (newImage *image.Gray) {
	b := img.Bounds()
	newImage = image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			newImage.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return
}

// DetectEdge for edge detection
func DetectEdge(img *image.Gray) *image.Gray {
	b := img.Bounds()
	//img := ImageToGray(colorImg)

	offset := 1
	newImage := image.NewGray(b)
	wg := sync.WaitGroup{}

	for y := b.Min.Y + offset; y < b.Max.Y-offset; y++ {
		y := y
		x1 := b.Min.X + offset
		x2 := b.Max.X - offset
		wg.Add(1)
		go doARow(y, x1, x2, img, newImage, &wg)
	}
	wg.Wait()
	return newImage
}
