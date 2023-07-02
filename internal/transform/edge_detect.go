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
	lumThresh := custom.NewNumericalEntry()
	lumThresh.SetText(fmt.Sprintf("%d", cfg.LuminanceThreshold))
	tryButton := widget.NewButton("Test", func() {})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "F", Widget: f},
			{Text: "S", Widget: s},
			{Text: "Luminance Threshold", Widget: lumThresh},
			{Widget: tryButton},
		},
		OnSubmit: func() {
			if xf, err := strconv.ParseFloat(f.Text, 64); err == nil {
				cfg.F = xf
			}
			if xs, err := strconv.ParseFloat(s.Text, 64); err == nil {
				cfg.S = xs
			}
			if lt, err := strconv.Atoi(lumThresh.Text); err == nil {
				cfg.LuminanceThreshold = uint8(lt)
			}
			globals.SetImage(transform.DetectEdge(c, cfg))
			globals.LoadImage()
			dialog.Close()
		},
		OnCancel: func() {
			dialog.Close()
		},
	}
	preview := container.New(
		layout.NewGridLayoutWithColumns(3),
		widget.NewLabel("Preview (Test)"),
		fi,
		widget.NewLabel(""))
	cfgForm := container.New(
		layout.NewGridLayoutWithColumns(3),
		container.New(
			layout.NewGridLayoutWithRows(10),
			widget.NewLabel("Configure F and S Parameters To Contol Edge Detection"),
			widget.NewLabel(""),
			widget.NewLabel("Set Luminance Threshold to Suppress Noise / Despeckle (0-255)")),
		form,
		container.New(
			layout.NewGridLayoutWithRows(10),
			widget.NewLabel("F and S control a sigmoid function over Edge Luminance such that:"),
			widget.NewLabel("Where Gradient is the convolved edge intensity, ranging from 0 to 256:"),
			widget.NewLabel("Edge Luminance = 255 * (1 + e ^ (-1.0*(Gradient-F)/S))"),
		),
	)
	layout := container.New(layout.NewGridLayoutWithRows(2), preview, cfgForm)
	tryButton.OnTapped = func() {
		if xf, err := strconv.ParseFloat(f.Text, 64); err == nil {
			cfg.F = xf
		}
		if xs, err := strconv.ParseFloat(s.Text, 64); err == nil {
			cfg.S = xs
		}
		if lt, err := strconv.Atoi(lumThresh.Text); err == nil {
			cfg.LuminanceThreshold = uint8(lt)
		}
		fi := canvas.NewImageFromImage(transform.DetectEdge(c, cfg))
		fi.FillMode = canvas.ImageFill(canvas.ImageFillContain)
		preview.Objects[1] = fi
	}

	dialog.SetContent(layout)
	dialog.Show()
}
