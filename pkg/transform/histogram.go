package transform

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/StephaneBunel/bresenham"
)

type Stats struct {
	counts        map[uint8]int
	MaxCount      uint8
	HighestLumins uint8
}

func NewStats() (s *Stats) {
	return &Stats{
		counts:        map[uint8]int{},
		MaxCount:      0,
		HighestLumins: 0,
	}
}
func (s *Stats) Bump(y uint8) {
	if y == 0 {
		return
	}
	s.counts[y] += 1
	if s.counts[y] > s.counts[s.MaxCount] {
		s.MaxCount = y
	}
	if y > s.HighestLumins {
		s.HighestLumins = y
	}
}

func (s *Stats) Histogram(cfg *EdgeDetectConfig) (img *image.RGBA) {
	barWidth := 5
	maxHeight := 1280.0
	heightScale := maxHeight / math.Log(float64(s.counts[s.MaxCount]))
	//fmt.Printf("heightScale: %f\n", heightScale)
	newBounds := image.Rectangle{image.Point{0, 0}, image.Point{X: 256 * barWidth, Y: int(math.Log(float64(s.counts[s.MaxCount]))*heightScale) + 5}}
	img = image.NewRGBA(newBounds)
	draw.Draw(img, newBounds, &image.Uniform{color.White}, image.Point{0, 0}, draw.Src)
	sig := (maxHeight / (1 + math.Exp(cfg.F/cfg.S)))
	leftPoint := image.Point{X: 0, Y: int(maxHeight-sig) + 1}
	for c := 0; c < 256; c += 1 {
		sig = (maxHeight / (1 + math.Exp(-1*(float64(c)-cfg.F)/cfg.S)))
		rightPoint := image.Point{X: c*barWidth + barWidth, Y: int(maxHeight-sig) + 1}

		// Histogram uses log y scale
		if s.counts[uint8(c)] > 0 {
			ul := image.Point{X: c * barWidth, Y: int(maxHeight - (math.Log(float64(s.counts[uint8(c)])) * heightScale))}
			lr := image.Point{X: c*barWidth + barWidth, Y: int(maxHeight)}
			r := image.Rectangle{ul, lr}
			draw.Draw(img, r, &image.Uniform{color.RGBA{0, 127, 255, 255}}, ul, draw.Src)
		}
		bresenham.DrawLine(img, leftPoint.X, leftPoint.Y, rightPoint.X, rightPoint.Y, color.RGBA{255, 0, 0, 255})
		bresenham.DrawLine(img, leftPoint.X+1, leftPoint.Y+1, rightPoint.X+1, rightPoint.Y+1, color.RGBA{255, 0, 0, 255})
		leftPoint = rightPoint
	}
	maxLog := (int(math.Log10(float64(s.counts[s.MaxCount]))))
	//fmt.Println("maxLog:", maxLog)
	for p := 0; p <= maxLog; p++ {
		height := float64(math.Pow10(p))
		bresenham.DrawLine(img, 0, int(maxHeight-math.Log(height)*heightScale), 256*barWidth, int(maxHeight-math.Log(height)*heightScale), color.Black)
	}

	return
}

func GrayScaleStats(img *image.Gray) (stats *Stats) {
	stats = NewStats()
	bounds := img.Bounds()
	for i := bounds.Min.Y; i <= bounds.Max.Y; i++ {
		for j := bounds.Min.X; j <= bounds.Max.X; j++ {
			stats.Bump(img.GrayAt(j, i).Y)
		}
	}
	return
}

func GetHist(img *image.Gray, cfg *EdgeDetectConfig) (newImg *image.RGBA) {
	stats := GrayScaleStats(img)
	//fmt.Printf("lumins with highest count: %d @ %d pixels.  Higehest Lumins: %d @ %d \n", stats.MaxCount, stats.counts[stats.MaxCount], stats.HighestLumins, stats.counts[stats.HighestLumins])
	newImg = stats.Histogram(cfg)
	return
}
