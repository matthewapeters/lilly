package layers

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/matthewapeters/lilly/internal/globals"
)

type Layer struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Order   int             `json:"order"`
	Bounds  image.Rectangle `json:"bounds"`
	Display bool            `json:"display"`
	Image   image.Image     `json:"image_b64"`
	Anchor  image.Point     `json:"anchor"`
}

var Layers = map[int]*Layer{}

func BaseLayer() {
	baseLayer := Layer{
		ID:      "baseLayer",
		Name:    globals.AppCtx.Value(globals.FileName).(string),
		Order:   0,
		Display: true,
		Image:   globals.GetImage(),
		Anchor:  globals.GetImage().Bounds().Min,
	}
	Layers[0] = &baseLayer
}

func Show() {
	layersWindow := globals.App.NewWindow("Layers")
	s := len(Layers)
	var layersGrid *fyne.Container
	objects := []fyne.CanvasObject{}
	for _, l := range Layers {
		objects = append(objects, l.GetWidget())
	}
	layersGrid = container.NewGridWithRows(s, objects...)
	layersWindow.SetContent(layersGrid)
	layersWindow.Show()
}

func Hide() {}

func Add() {}

func ReOrder() {}

func (l *Layer) GetWidget() (layerContainer *fyne.Container) {
	check := widget.NewCheck(l.Name, func(changed bool) {})
	entry := widget.NewEntry()
	entry.Text = l.Name
	entry.OnChanged = func(newName string) {
		l.Name = newName
		check.Text = newName
	}
	layerContainer = container.NewGridWithColumns(3,
		check,
		entry,
	)
	return
}
func LoadImage() {
	if canvs := globals.GetWindow(); canvs != nil {
		canvs.SetTitle(fmt.Sprintf("%s", globals.AppCtx.Value(globals.FilePath)))
		if img := globals.GetImage(); img != nil {
			i := canvas.NewImageFromImage(img)
			i.FillMode = canvas.ImageFill(canvas.ImageFillContain)
			canvs.SetContent(i)
		}
	}
}
