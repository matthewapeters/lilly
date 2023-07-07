package globals

import (
	"context"
	"image"

	"fyne.io/fyne/v2"
)

type WindowType string
type ImageType string
type FileNameType string
type FilePathType string
type URIType string

const (
	FileName  = FileNameType("FileName")
	FilePath  = FilePathType("FilePath")
	FileURI   = URIType("FileURI")
	AppWindow = WindowType("AppWindow")
	AppImage  = ImageType("AppImage")
)

var (
	AppCtx = context.Context(context.Background())
	App    fyne.App
)

func GetImage() image.Image {
	imgRaw := AppCtx.Value(AppImage)
	if imgRaw == nil {
		return nil
	}
	return imgRaw.(image.Image)
}

func SetImage(img image.Image) {
	AppCtx = context.WithValue(AppCtx, AppImage, img)
}

func GetWindow() fyne.Window {
	wndRaw := AppCtx.Value(AppWindow)
	if wndRaw == nil {
		return nil
	}
	return wndRaw.(fyne.Window)
}

func SetValue(key any, value any) {
	AppCtx = context.WithValue(AppCtx, key, value)
}

func GetLastFile() fyne.URI {
	u := AppCtx.Value(FileURI)
	if u == nil {
		return nil
	}
	return u.(fyne.URI)
}

var InfoMenu = fyne.NewMenuItem(
	"File Info",
	func() {})

var TransfomEdgeDetect = fyne.NewMenuItem(
	"Edge Detect",
	func() {})

var TransformScale = fyne.NewMenuItem(
	"Scale",
	func() {})

var FileSaveAs = fyne.NewMenuItem(
	"Save As",
	func() {})
