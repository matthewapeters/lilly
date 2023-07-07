package main

import (
	"context"
	"fmt"

	"github.com/matthewapeters/lilly/internal/globals"
	"github.com/matthewapeters/lilly/internal/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

const (
	ApplicationName = "Lilly Image Utility"
)

var (
	ctxChan = make(chan context.Context)
)

func updateCtx() {
	for {
		// Read from the ctxChan until it is closed in tidyUp
		newCtx := <-ctxChan
		//fmt.Println("updated context")
		if newCtx == nil {
			return
		}
		globals.AppCtx = newCtx
	}
}

func tidyUp() {
	fmt.Println("Exited")
	close(ctxChan)
}

func main() {
	globals.App = app.NewWithID(ApplicationName)
	globals.App.SetIcon(theme.ColorPaletteIcon())
	myWindow := globals.App.NewWindow(ApplicationName)
	myWindow.Resize(fyne.NewSize(400, 400))
	globals.AppCtx = context.WithValue(globals.AppCtx, globals.AppWindow, myWindow)
	globals.AppCtx = context.WithValue(globals.AppCtx, globals.FileURI, nil)
	menu.InitialMenuLoad(ctxChan)
	go updateCtx()
	myWindow.Show()
	globals.App.Run()
	tidyUp()
}
