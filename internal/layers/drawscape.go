package layers

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DrawScape struct {
	widget.BaseWidget
	Img       *canvas.Image
	Bounds    image.Rectangle
	Selection image.Rectangle
}

type drawScapeRenderer struct {
	layout  fyne.Layout
	objects []fyne.CanvasObject
}

func NewDrawScape(img image.Image) *DrawScape {
	i := canvas.NewImageFromImage(img)
	i.FillMode = canvas.ImageFill(canvas.ImageFillContain)
	ds := &DrawScape{
		Img:    i,
		Bounds: img.Bounds(),
	}
	ds.ExtendBaseWidget(ds)
	return ds
}

func (ds *DrawScape) CreateRenderer() fyne.WidgetRenderer {
	return &drawScapeRenderer{
		layout:  layout.NewMaxLayout(),
		objects: []fyne.CanvasObject{ds.Img},
	}
}

func (dsr *drawScapeRenderer) MinSize() fyne.Size {
	return dsr.layout.MinSize(dsr.objects)
}

func (dsr *drawScapeRenderer) Layout(s fyne.Size) {
	dsr.layout.Layout(dsr.objects, s)
}

func (dsr *drawScapeRenderer) Refresh() {
	canvas.Refresh(dsr.objects[0])
}

func (dsr *drawScapeRenderer) Objects() []fyne.CanvasObject {
	return dsr.objects
}

func (dsr *drawScapeRenderer) Destroy() {
}

func (ds *DrawScape) Tapped(pe *fyne.PointEvent) {
	fmt.Println(pe.Position, ds.Img.Visible(), ds.MinSize())
}

func (ds *DrawScape) TappedSecondary(pe *fyne.PointEvent) {
	fmt.Println("Secondary", pe.Position, ds.Img.Visible(), ds.MinSize())
}

func (ds *DrawScape) GetSelection() {

}
