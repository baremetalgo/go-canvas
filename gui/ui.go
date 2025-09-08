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
	ui.LayerEditor = widgets.NewLayerEditor("Layer Editor")
	ui.ToolBox = widgets.NewToolBox("Tools", ui)
	ui.ColorSet = widgets.NewColorSet("Colors", ui.ToolBox)
	ui.PaintCanvas = widgets.NewPaintCanvas("Canvas", ui.LayerEditor, ui.ColorSet)

	ui.MenuBar = widgets.NewMenuBar("Menu")
	ui.BrushPattenBar = widgets.NewBrushPatternBar("Patterns", ui.PaintCanvas)
	ui.StatusBar = widgets.NewStatusBar("StatusBar")

	return ui
}

func (u *UserInterface) GetLayerEditor() *widgets.LayerEditor {
	return u.LayerEditor
}

func (u *UserInterface) GetCanvas() *widgets.PaintCanvas {
	return u.PaintCanvas
}

func (u *UserInterface) Draw() {
	u.PaintCanvas.Draw()
	u.MenuBar.Draw()
	u.ToolBox.Draw()
	u.ColorSet.Draw()
	u.BrushPattenBar.Draw()
	u.LayerEditor.Draw()
	u.StatusBar.Draw()

}
