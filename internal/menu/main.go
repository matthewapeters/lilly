package menu

import (
	"context"

	"fyne.io/fyne/v2"
	"github.com/matthewapeters/fyne_stuff/internal/globals"
	"github.com/matthewapeters/fyne_stuff/internal/lillyfile"
)

func InitialMenuLoad(ctxChan chan context.Context) {
	canvs := globals.AppCtx.Value(globals.AppWindow).(fyne.Window)
	globals.InfoMenu.Disabled = true
	globals.TransformScale.Disabled = true
	globals.TransfomEdgeDetect.Disabled = true
	globals.FileSaveAs.Disabled = true

	fileMenu := fyne.NewMenu(
		"File",
		fyne.NewMenuItem(
			"Open",
			func() { lillyfile.OpenFile() }),
		globals.FileSaveAs,
		globals.InfoMenu,
	)

	menu := fyne.NewMainMenu(
		fileMenu,
		fyne.NewMenu(
			"Transform",
			globals.TransfomEdgeDetect,
			globals.TransformScale),
	)

	canvs.SetMainMenu(menu)
}
