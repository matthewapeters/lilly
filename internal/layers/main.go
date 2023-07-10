package layers

import (
	"fmt"
	"image"
	"image/draw"
	"sort"

	"fyne.io/fyne/v2"
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

func Add() {
	globals.SetValue(globals.DrawSelection, true)
}

func ReOrder() {}

func (l *Layer) GetWidget() (layerContainer *fyne.Container) {
	check := widget.NewCheck(l.Name, func(changed bool) {})
	check.SetChecked(l.Display)
	check.OnChanged = func(tf bool) {
		l.Display = tf
		LoadImage()
	}
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

func applyLayers(img *image.RGBA) {
	layerIdxs := []int{}
	for key := range Layers {
		if Layers[key].Display {
			layerIdxs = append(layerIdxs, key)
		}
	}
	sort.IntSlice(layerIdxs).Sort()
	for _, key := range layerIdxs {
		i := Layers[key].Image
		p := Layers[key].Anchor
		draw.Draw(img, i.Bounds(), i, p, draw.Src)
	}
}

func LoadImage() {
	if win := globals.GetWindow(); win != nil {
		win.SetTitle(fmt.Sprintf("%s", globals.AppCtx.Value(globals.FilePath)))
		if img := globals.GetImage(); img != nil {
			base := image.NewRGBA(globals.AppCtx.Value(globals.Bounds).(image.Rectangle))
			applyLayers(base)
			ds := NewDrawScape(base)
			globals.SetValue(globals.AppDrawScape, ds)
			win.SetContent(ds)
			win.Show()
		}
	}
}
