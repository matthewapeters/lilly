package transform

import (
	"fmt"
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/matthewapeters/lilly/internal/custom"
	"github.com/matthewapeters/lilly/internal/globals"
	"github.com/matthewapeters/lilly/pkg/transform"
)

func EdgeDetect() {
	cfg := transform.DefaultEdgeDetectConfig()
	i := globals.GetImage()
	c, ok := i.(*image.RGBA)
	if !ok {
		fmt.Println("could not treat image as RGBA")
		return
	}

	dialog := fyne.CurrentApp().NewWindow("Edge Detection")
	dialog.SetFixedSize(false)
	dialog.Resize(fyne.NewSize(840, 840))

	fi := canvas.NewImageFromImage(i)
	fi.FillMode = canvas.ImageFill(canvas.ImageFillContain)
	f := custom.NewNumericalEntry()
	f.SetText(fmt.Sprintf("%f", cfg.F))
	s := custom.NewNumericalEntry()
	s.SetText(fmt.Sprintf("%f", cfg.S))
	tryButton := widget.NewButton("Test", func() {})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "F", Widget: f},
			{Text: "S", Widget: s},
			{Widget: tryButton},
		},
		OnSubmit: func() {
			globals.SetImage(transform.DetectEdge(c, cfg))
			globals.LoadImage()
			dialog.Close()
		},
		OnCancel: func() {
			dialog.Close()
		},
	}
	layout := container.New(layout.NewGridLayoutWithRows(2), form, fi)
	tryButton.OnTapped = func() {
		if xf, err := strconv.ParseFloat(f.Text, 64); err == nil {
			cfg.F = xf
		}
		if xs, err := strconv.ParseFloat(s.Text, 64); err == nil {
			cfg.S = xs
		}
		fi := canvas.NewImageFromImage(transform.DetectEdge(c, cfg))
		fi.FillMode = canvas.ImageFill(canvas.ImageFillContain)
		layout.Objects[1] = fi
	}

	dialog.SetContent(layout)
	dialog.Show()
}
