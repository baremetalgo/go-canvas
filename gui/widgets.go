package widgets

type Widgets struct {
	ToolBox        *ToolBox
	MenuBar        *MenuBar
	ColorSet       *ColorSet
	BrushPattenBar *BrushPatternBar
	LayerEditor    *LayerEditor
	StatusBar      *StatusBar
}

func NewWidgets() *Widgets {
	widgets := &Widgets{}
	widgets.ToolBox = NewToolBox("Tools")
	widgets.MenuBar = NewMenuBar("Menu")
	widgets.ColorSet = NewColorSet("Colors", widgets.ToolBox)
	widgets.BrushPattenBar = NewBrushPatternBar("Patterns")
	widgets.LayerEditor = NewLayerEditor("Layers")
	widgets.StatusBar = NewStatusBar("StatusBar")
	return widgets
}

func (w *Widgets) Draw() {
	w.MenuBar.Draw()
	w.ToolBox.Draw()
	w.ColorSet.Draw()
	w.BrushPattenBar.Draw()
	w.LayerEditor.Draw()
	w.StatusBar.Draw()

}
