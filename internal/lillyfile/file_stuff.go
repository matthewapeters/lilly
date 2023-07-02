package lillyfile

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/matthewapeters/lilly/internal/globals"
	"github.com/matthewapeters/lilly/internal/transform"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func OpenFile() {
	openFileCallBack := func(readCloser fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		// No file was selected (Cancel button) just return
		if readCloser == nil {
			return
		}
		globals.SetValue(globals.FileName, readCloser.URI().Name())
		globals.SetValue(globals.FilePath, readCloser.URI().Path())
		parent := "file:/"
		parts := strings.Split(readCloser.URI().Path(), "/")
		for i := 0; i < len(parts)-1; i++ {
			parent += fmt.Sprintf("/%s", parts[i])
		}
		//fmt.Printf("Debug: parent: %s\n", parent)
		parentURI, err := storage.ParseURI(parent)
		if err != nil {
			fmt.Println(err)
			os.Exit(7)
		}
		globals.SetValue(globals.FileURI, parentURI)
		img, _, err := image.Decode(readCloser)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		globals.SetImage(img)
		globals.InfoMenu.Disabled = false
		globals.InfoMenu.Action = FileInfo
		globals.TransfomEdgeDetect.Disabled = false
		globals.TransfomEdgeDetect.Action = transform.EdgeDetect
		globals.TransformScale.Disabled = false
		globals.FileSaveAs.Disabled = false
		globals.FileSaveAs.Action = SaveFileAs
		globals.LoadImage()
	}

	picker := dialog.NewFileOpen(openFileCallBack, globals.GetWindow())
	if uri := globals.GetLastFile(); uri != nil {
		listable, err := storage.ListerForURI(uri)
		if err != nil {
			fmt.Println(err, uri.Path())
			os.Exit(6)
		}
		picker.SetLocation(listable)
	}
	picker.Show()
}

func SaveFileAs() {
	saveFileCallback := func(writeCloser fyne.URIWriteCloser, err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		// No file was selected (Cancel button) just return
		if writeCloser == nil {
			return
		}
		if img := globals.GetImage(); img != nil {
			err = png.Encode(writeCloser, img)
			if err != nil {
				fmt.Println(err)
				os.Exit(8)
			}
		}
	}

	picker := dialog.NewFileSave(saveFileCallback, globals.GetWindow())
	if uri := globals.GetLastFile(); uri != nil {
		listable, err := storage.ListerForURI(uri)
		if err != nil {
			fmt.Println(err, uri.Path())
			os.Exit(6)
		}
		picker.SetLocation(listable)
	}
	picker.Show()
}

func FileInfo() {
	infoWindow := fyne.CurrentApp().NewWindow(fmt.Sprintf("Info: %s", globals.GetLastFile().Path()))
	pathWidget := widget.NewLabel(
		fmt.Sprintf(
			"Path: %s",
			globals.AppCtx.Value(globals.FilePath)))
	pathWidget.Wrapping = fyne.TextWrapWord
	dimensionWidget := widget.NewLabel(
		fmt.Sprintf(
			"Dimensions: %d x %d",
			globals.GetImage().Bounds().Max.X-globals.GetImage().Bounds().Min.X,
			globals.GetImage().Bounds().Max.Y-globals.GetImage().Bounds().Min.Y),
	)
	dimensionWidget.Wrapping = fyne.TextWrapOff

	infoObject := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(840, 50)),
		pathWidget,
		dimensionWidget,
	)

	infoWindow.SetContent(infoObject)
	infoWindow.Show()
}
