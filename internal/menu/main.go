package menu

import (
	"context"

	"fyne.io/fyne/v2"
	"github.com/matthewapeters/lilly/internal/globals"
	"github.com/matthewapeters/lilly/internal/layers"
	"github.com/matthewapeters/lilly/internal/lillyfile"
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

	transformMenu := fyne.NewMenu(
		"Transform",
		globals.TransfomEdgeDetect,
		globals.TransformScale)

	layersMenu := fyne.NewMenu(
		"Layers",
		fyne.NewMenuItem("Show Layers", layers.Show),
	)

	menu := fyne.NewMainMenu(
		fileMenu,
		transformMenu,
		layersMenu,
	)

	canvs.SetMainMenu(menu)
}
