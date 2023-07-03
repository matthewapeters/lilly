package transform

import (
	"fmt"
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/matthewapeters/lilly/internal/globals"
	"github.com/matthewapeters/lilly/pkg/transform"
)

type DataBinder struct {
	binder binding.Float
}

func (db DataBinder) Get() (string, error) {
	v, err := db.binder.Get()
	return fmt.Sprintf("%f", v), err
}

func (db DataBinder) AddListener(bs binding.DataListener) {
	db.binder.AddListener(bs)
}
func (db DataBinder) RemoveListener(bs binding.DataListener) {
	db.binder.RemoveListener(bs)
}
func (db DataBinder) Set(v string) error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	db.binder.Set(f)
	return nil
}

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

	tData := binding.NewFloat()
	tData.Set(cfg.F)
	tDb := DataBinder{binder: tData}
	tSlider := widget.NewSliderWithData(1, 1020, tData)
	tLabel := widget.NewLabel("")
	tLabel.Bind(tDb)

	sData := binding.NewFloat()
	sData.Set(cfg.S)
	sDb := DataBinder{binder: sData}
	sliderS := widget.NewSliderWithData(0.1, 2040, sData)
	sliderLabel := widget.NewLabel("")
	sliderLabel.Bind(sDb)
	//lumThresh := custom.NewNumericalEntry()
	//lumThresh.SetText(fmt.Sprintf("%d", cfg.LuminanceThreshold))
	tryButton := widget.NewButton("Test", func() {})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "T", Widget: tSlider},
			{Text: "", Widget: tLabel},
			{Text: "S", Widget: sliderS},
			{Text: "", Widget: sliderLabel},
			//{Text: "Luminance Threshold", Widget: lumThresh},
			{Widget: tryButton},
		},
		OnSubmit: func() {
			cfg.F, _ = tData.Get()
			cfg.S, _ = sData.Get()
			//if lt, err := strconv.Atoi(lumThresh.Text); err == nil {
			//	cfg.LuminanceThreshold = uint8(lt)
			//}
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
			widget.NewLabel("Configure T and S Parameters To Contol Edge Detection"),
			widget.NewLabel(""),
			//widget.NewLabel("Set Luminance Threshold to Suppress Noise / Despeckle (0-255)")
		),
		form,
		container.New(
			layout.NewGridLayoutWithRows(10),
			widget.NewLabel("T and S control a sigmoid function over Edge Luminance such that:"),
			widget.NewLabel("Where Gradient is the convolved edge intensity, ranging from 0 to 1020:"),
			widget.NewLabel("Where T is the Threshold, mapping to the Sigmoid inflection point X value"),
			widget.NewLabel("Where S controls the tangent at the Sigmoid inflection point"),
			widget.NewLabel("NOTE: Reducing S to 0.1 is effectively a Step Function; 64 results in 45 degree sigmoid"),
			widget.NewLabel("Edge Luminance = 255 / (1 + e ^ (-1.0*(Gradient-T)/S))"),
		),
	)
	layout := container.New(layout.NewGridLayoutWithRows(2), preview, cfgForm)
	tryButton.OnTapped = func() {
		cfg.F, _ = tData.Get()
		cfg.S, _ = sData.Get()
		//if lt, err := strconv.Atoi(lumThresh.Text); err == nil {
		//	cfg.LuminanceThreshold = uint8(lt)
		//}
		fi := canvas.NewImageFromImage(transform.DetectEdge(c, cfg))
		fi.FillMode = canvas.ImageFill(canvas.ImageFillContain)
		preview.Objects[1] = fi
	}

	dialog.SetContent(layout)
	dialog.Show()
}
