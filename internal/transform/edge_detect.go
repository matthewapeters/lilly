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

	// Display the loaded image as Grayscale
	bw := transform.ImageToGray(c)
	bwi := canvas.NewImageFromImage(bw)
	bwi.FillMode = canvas.ImageFill(canvas.ImageFillContain)

	// Display the loaded image after edges are detected
	edges := transform.DetectEdge(bw)
	edgesImg := canvas.NewImageFromImage(edges)
	edgesImg.FillMode = canvas.ImageFill(canvas.ImageFillContain)

	hist := transform.GetHist(edges, cfg)
	histImg := canvas.NewImageFromImage(hist)
	histImg.FillMode = canvas.ImageFill(canvas.ImageFillContain)

	final := transform.ApplySigmoid(edges, cfg)
	finalImg := canvas.NewImageFromImage(final)
	finalImg.FillMode = canvas.ImageFill(canvas.ImageFillContain)

	tData := binding.NewFloat()
	tData.Set(cfg.F)
	tDb := DataBinder{binder: tData}
	tSlider := widget.NewSliderWithData(0, 255, tData)
	tLabel := widget.NewLabel("")
	tLabel.Bind(tDb)

	sData := binding.NewFloat()
	sData.Set(cfg.S)
	sDb := DataBinder{binder: sData}
	sliderS := widget.NewSliderWithData(0.1, 65, sData)
	sliderLabel := widget.NewLabel("")
	sliderLabel.Bind(sDb)
	tryButton := widget.NewButton("Test Edge Enhance", func() {})
	tryButtonContainer := container.New(
		layout.NewGridLayoutWithColumns(3),
		tryButton,
	)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "T", Widget: tSlider},
			{Text: "", Widget: tLabel},
			{Text: "S", Widget: sliderS},
			{Text: "", Widget: sliderLabel},
			{Widget: tryButtonContainer},
		},
		OnSubmit: func() {
			cfg.F, _ = tData.Get()
			cfg.S, _ = sData.Get()
			globals.SetImage(transform.ApplySigmoid(edges, cfg))
			globals.LoadImage()
			dialog.Close()
		},
		OnCancel: func() {
			dialog.Close()
		},
	}

	// Give the preview container a variable so we can
	// update the content when Test button is pressed
	preview := container.New(
		layout.NewMaxLayout(),
		finalImg)

	grayContainer := container.NewBorder(
		widget.NewLabel("Loaded Image as Grayscale"),
		nil, nil, nil,
		bwi,
	)

	edgesContainer := container.NewBorder(
		widget.NewLabel("Default Edge Detection"),
		nil, nil, nil,
		edgesImg)

	previewContainer := container.NewBorder(
		widget.NewLabel("Enhanced Edge Detection (Final)"),
		nil, nil, nil,
		preview)

	images := container.New(
		layout.NewGridLayoutWithColumns(3),
		grayContainer,
		edgesContainer,
		previewContainer)

	histImgContainer := container.NewBorder(
		widget.NewLabel("Histogram of Edge Luminosity - Y-axis is Log10 Scale"),
		nil, nil, nil,
		histImg,
	)

	cfgForm := container.New(
		layout.NewGridLayoutWithColumns(3),
		container.New(
			layout.NewGridLayoutWithRows(10),
			widget.NewLabel("Configure T and S Parameters To Enhance Edge Brightness"),
			widget.NewLabel("T and S control a sigmoid function over Edge Luminance such that:"),
			widget.NewLabel("Where Gradient is the convolved edge intensity, ranging from 0 to 255:"),
			widget.NewLabel("Where T is the Threshold, mapping to the Sigmoid inflection point X value"),
			widget.NewLabel("Where S controls the tangent at the Sigmoid inflection point"),
			widget.NewLabel("NOTE: Reducing S to 0.1 is effectively a Step Function; 64 results in 45 degree sigmoid"),
			widget.NewLabel("Edge Luminance = 255 / (1 + e ^ (-1.0*(Gradient-T)/S))"),
		),
		form,
		histImgContainer,
	)
	layout := container.New(layout.NewGridLayoutWithRows(2),
		images,
		cfgForm)
	tryButton.OnTapped = func() {
		cfg.F, _ = tData.Get()
		cfg.S, _ = sData.Get()
		finalImg = canvas.NewImageFromImage(transform.ApplySigmoid(edges, cfg))
		finalImg.FillMode = canvas.ImageFill(canvas.ImageFillContain)
		preview.Objects[0] = finalImg
		hist := transform.GetHist(edges, cfg)
		histImg := canvas.NewImageFromImage(hist)
		histImg.FillMode = canvas.ImageFill(canvas.ImageFillContain)
		histImgContainer.Objects[0] = histImg
	}

	dialog.SetContent(layout)
	dialog.Show()
}
