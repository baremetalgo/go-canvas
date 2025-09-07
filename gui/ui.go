package gui

import (
	"go-canvas/widgets"
)

type UserInterface struct {
	ToolBox        *widgets.ToolBox
	MenuBar        *widgets.MenuBar
	ColorSet       *widgets.ColorSet
	BrushPattenBar *widgets.BrushPatternBar
	LayerEditor    *widgets.LayerEditor
	StatusBar      *widgets.StatusBar
	PaintCanvas    *widgets.PaintCanvas
}

func NewUserInterface() *UserInterface {
	ui := &UserInterface{}
	ui.ToolBox = widgets.NewToolBox("Tools")
	ui.MenuBar = widgets.NewMenuBar("Menu")
	ui.ColorSet = widgets.NewColorSet("Colors", ui.ToolBox)
	ui.BrushPattenBar = widgets.NewBrushPatternBar("Patterns")
	ui.LayerEditor = widgets.NewLayerEditor("Layer Editor")
	ui.StatusBar = widgets.NewStatusBar("StatusBar")
	ui.PaintCanvas = widgets.NewPaintCanvas("Canvas", ui.LayerEditor)
	return ui
}

func (w *UserInterface) Draw() {
	w.PaintCanvas.Draw()
	w.MenuBar.Draw()
	w.ToolBox.Draw()
	w.ColorSet.Draw()
	w.BrushPattenBar.Draw()
	w.LayerEditor.Draw()
	w.StatusBar.Draw()

}
