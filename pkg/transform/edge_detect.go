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

	//
	F         float64
	S         float64
	ShowAngle bool
	//LuminanceThreshold uint8
}

func DefaultEdgeDetectConfig() *EdgeDetectConfig {
	return &EdgeDetectConfig{
		RedFactor:   0.299,
		GreenFactor: 0.587,
		BlueFactor:  0.114,
		F:           127,
		S:           5.0,
		ShowAngle:   false,
		//LuminanceThreshold: 127,
	}
}

// Luminance Convert color to lumins
func (cfg *EdgeDetectConfig) Luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (float64(r) * cfg.RedFactor) + (float64(g) * cfg.GreenFactor) + (float64(b) * cfg.BlueFactor)
}

func computePixelGray(img *image.Gray, x, y int, cfg *EdgeDetectConfig) (float64, float64) {
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

func doARow(y int, x1, x2 int, img *image.Gray, newimg *image.RGBA, cfg *EdgeDetectConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	var o, r, g, b uint8
	// The convolved value of a single point
	var gradient float64
	var angle float64
	for x := x1; x < x2; x++ {
		// compute gradient and ganle
		gradient, angle = computePixelGray(img, x, y, cfg)
		//fmt.Println(gradient)

		// apply sigmoid to enhance dominant gradients
		o = uint8(255 / (1 + math.Exp(-1.0*(gradient-cfg.F)/cfg.S)))
		r = o
		g = r
		b = r
		if o > 200 && cfg.ShowAngle {
			//a2 is the angle in radians scaled to -255 to 255
			a2 := int((angle * 255.00 / math.Pi) * 100.00)

			// White-to-Red quadrant
			if math.Pi/-2.0 >= angle && angle > math.Pi*-1.0 {
				r = 255
				g = 255 - uint8(a2)*2
				b = 255 - uint8(a2)*2
			}
			// Yellow-to-green quadrant
			if 0 >= angle && angle > math.Pi/-2 {
				r = 255 - (uint8(127+a2) * 2)
				g = 255 - (255 - (uint8(127+a2) * 2))
				b = 0
			}
			//Green-to-blue quadrant
			if math.Pi/2 > angle && angle >= 0 {
				r = 0
				g = 255
				b = 255 - (255 - uint8(127+a2)*2)
			}
			// blue-to-magenta quadrant
			if angle >= math.Pi/2 {
				r = uint8(a2) * 2
				g = 0
				b = 255
			}
			newimg.Set(x, y, color.RGBA{r, g, b, o})
		} else {
			//if o > cfg.LuminanceThreshold {
			newimg.Set(x, y, color.Gray{o})
			//}
		}
	}
}

// ImageToGray converts an RGBA image to Gray
func ImageToGray(img *image.RGBA) *image.Gray {
	b := img.Bounds()
	newImage := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			newImage.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return newImage
}

// DetectEdge for edge detection
func DetectEdge(colorImg *image.RGBA, cfg *EdgeDetectConfig) *image.RGBA {
	b := colorImg.Bounds()
	img := ImageToGray(colorImg)

	offset := 1
	newImage := image.NewRGBA(b)
	wg := sync.WaitGroup{}

	for y := b.Min.Y + offset; y < b.Max.Y-offset; y++ {
		y := y
		x1 := b.Min.X + offset
		x2 := b.Max.X - offset
		wg.Add(1)
		go doARow(y, x1, x2, img, newImage, cfg, &wg)
	}
	wg.Wait()
	return newImage
}
