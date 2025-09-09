package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ToolBoxIcons map[string]rl.Texture2D

type GUI interface {
	GetLayerEditor() *LayerEditor
	GetCanvas() *PaintCanvas
}

type ToolBox struct {
	BaseWidget
	UserInterface     GUI
	CircleBrushButton rl.Rectangle
	SquareBrushButton rl.Rectangle
	ShapeDrawButton   rl.Rectangle
	EraserDrawButton  rl.Rectangle
	ClearDrawButton   rl.Rectangle
	CircleShapeButton rl.Rectangle
	SquareShapeButton rl.Rectangle
	LineShapeButton   rl.Rectangle
	BucketShapeButton rl.Rectangle

	SelectedButton  rl.Rectangle
	BrushColor      rl.Color
	HightlightColor rl.Color
	DebugDraw       bool
	BrushMap        map[rl.Texture2D]*rl.Rectangle
	button_order    [8]*rl.Rectangle
}

func NewToolBox(name string, ui GUI) *ToolBox {
	tool_box := ToolBox{}
	tool_box.LoadToolBoxIcons()
	tool_box.UserInterface = ui
	// Setup BaseWidget first
	tool_box.BrushColor = rl.White
	tool_box.Visible = true
	tool_box.Name = name
	tool_box.BodyColor = rl.NewColor(50, 50, 50, 20)
	tool_box.DrawTitleBar = true
	tool_box.TextColor = rl.White
	tool_box.TitleBarHeight = 20
	tool_box.TitleBarColor = rl.DarkGray
	tool_box.Width = 75
	tool_box.Height = 350

	tool_box.Bounds = rl.NewRectangle(10, 10, float32(tool_box.Width), float32(tool_box.Height))
	tool_box.TitleBarBounds = rl.NewRectangle(tool_box.Bounds.X, tool_box.Bounds.Y, tool_box.Bounds.Width, tool_box.Bounds.Height)
	tool_box.IsDragging = false
	tool_box.HightlightColor = rl.NewColor(255, 255, 255, 100)
	tool_box.DebugDraw = false

	// Create brush map
	tool_box.BrushMap = map[rl.Texture2D]*rl.Rectangle{
		ToolBoxIcons["line_brush"]:   &tool_box.LineShapeButton,
		ToolBoxIcons["square_brush"]: &tool_box.SquareBrushButton,
		ToolBoxIcons["round_brush"]:  &tool_box.CircleBrushButton,
		ToolBoxIcons["square_shape"]: &tool_box.SquareShapeButton,
		ToolBoxIcons["round_shape"]:  &tool_box.CircleShapeButton,
		ToolBoxIcons["bucket_brush"]: &tool_box.BucketShapeButton,
		ToolBoxIcons["clear"]:        &tool_box.ClearDrawButton,
		ToolBoxIcons["eraser"]:       &tool_box.EraserDrawButton,
	}
	tool_box.button_order = [8]*rl.Rectangle{&tool_box.SquareBrushButton, &tool_box.CircleBrushButton,
		&tool_box.LineShapeButton, &tool_box.SquareShapeButton,
		&tool_box.CircleShapeButton, &tool_box.BucketShapeButton,
		&tool_box.EraserDrawButton, &tool_box.ClearDrawButton}

	return &tool_box
}

func (t *ToolBox) Draw() {
	t.BaseWidget.Draw()
	t.Update()

	for text, button := range t.BrushMap {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(text.Width), Height: float32(text.Height)}
		dest := rl.NewRectangle(button.X, button.Y, button.Width, button.Height)
		rl.DrawTexturePro(
			text,
			src,
			dest,
			rl.Vector2{X: 0, Y: 0},
			0,
			t.BrushColor)
	}
}

func (t *ToolBox) Update() {
	prev_y := float32(0)
	tool_box_height := int32(0)
	mouse_pos := rl.GetMousePosition()

	for idx, button := range t.button_order {
		button.X = t.Bounds.X + 14
		button.Y = t.Bounds.Y + prev_y
		button.Width = 40
		button.Height = 40
		prev_y += button.Height + 5

		if idx == 0 {
			button.Y = t.Bounds.Y + float32(t.TitleBarHeight)
			prev_y += 15
		}
		tool_box_height += int32(button.Height)

		if rl.CheckCollisionPointRec(mouse_pos, *button) {
			rl.DrawRectangleLinesEx(*button, 10, t.HightlightColor)
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				t.SelectedButton = *button
				if button == &t.ClearDrawButton {
					t.UserInterface.GetLayerEditor().ActiveLayer.Clear()
				}
				if button == &t.CircleBrushButton {
					t.UserInterface.GetCanvas().Brush.Shape = Circle
				}
				if button == &t.SquareBrushButton {
					t.UserInterface.GetCanvas().Brush.Shape = Rectangle
				}
				if button == &t.CircleShapeButton {
					t.UserInterface.GetCanvas().Brush.Shape = Round
				}
				if button == &t.SquareShapeButton {
					t.UserInterface.GetCanvas().Brush.Shape = Square
				}
				if button == &t.EraserDrawButton {
					t.UserInterface.GetCanvas().Brush.Color = rl.Blank
					t.UserInterface.GetCanvas().ColorSet.CurrentColor = rl.Blank

				}
				if button == &t.LineShapeButton {
					t.UserInterface.GetCanvas().Brush.Shape = Line

				}
				if button == &t.BucketShapeButton {
					t.UserInterface.GetCanvas().Brush.Shape = Bucket

				}
				t.UserInterface.GetCanvas().Brush.UsePattern = false
			}

		}

	}

	// set overall height
	t.Height = tool_box_height + t.TitleBarHeight + 50
}

func (t *ToolBox) LoadToolBoxIcons() {
	ToolBoxIcons = map[string]rl.Texture2D{
		"clear":        rl.LoadTexture("icons/clear_icon.png"),
		"eraser":       rl.LoadTexture("icons/eraser_icon.png"),
		"round_brush":  rl.LoadTexture("icons/round_brush.png"),
		"square_brush": rl.LoadTexture("icons/square_brush.png"),
		"round_shape":  rl.LoadTexture("icons/circle_shape.png"),
		"square_shape": rl.LoadTexture("icons/square_shape.png"),
		"line_brush":   rl.LoadTexture("icons/line_icon.png"),
		"bucket_brush": rl.LoadTexture("icons/bucket_icon.png"),
	}
}
