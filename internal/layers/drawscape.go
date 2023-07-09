package layers

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type DrawScape struct {
	*fyne.Container
	Bounds image.Rectangle
}

func NewDrawScape(img image.Image) *DrawScape {
	i := canvas.NewImageFromImage(img)
	i.FillMode = canvas.ImageFill(canvas.ImageFillContain)
	ds := &DrawScape{
		container.New(
			layout.NewMaxLayout(),
			i),
		img.Bounds(),
	}
	return ds
}

func (ds *DrawScape) Tapped(pe *fyne.PointEvent) {
	fmt.Println(pe.Position, ds.Objects[0].Visible(), ds.MinSize())
}

func (ds *DrawScape) TappedSecondary(pe *fyne.PointEvent) {
	fmt.Println("Secondary", pe.Position, ds.Objects[0].Visible(), ds.MinSize())
}
