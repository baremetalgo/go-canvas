package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Default_Widget_Header_Font rl.Font
var Default_tiny_Font rl.Font

func InitializeFonts() {
	Default_Widget_Header_Font = rl.LoadFontEx("fonts/CALIBRI.TTF", 14, nil, 0)
	Default_tiny_Font = rl.LoadFontEx("fonts/CALIBRI.TTF", 10, nil, 0)
}

type BaseWidget struct {
	Name            string
	Visible         bool
	Bounds          rl.Rectangle
	TitleBarBounds  rl.Rectangle
	TextColor       rl.Color
	BorderColor     rl.Color
	BodyColor       rl.Color
	DrawTitleBar    bool
	TitleBarHeight  int32
	TitleBarColor   rl.Color
	Width           int32
	Height          int32
	BorderThickness float32
	DrawBorder      bool

	IsDragging  bool
	DragOffsetX float32
	DragOffsetY float32
}

func (b *BaseWidget) Draw() {
	if !b.Visible {
		return
	}
	b.Update()

	// Draw body
	rl.DrawRectangle(b.Bounds.ToInt32().X, b.Bounds.ToInt32().Y, b.Width, b.Height, b.BodyColor)

	// Draw titleBar
	if b.DrawTitleBar {
		rl.DrawRectangle(b.TitleBarBounds.ToInt32().X, b.TitleBarBounds.ToInt32().Y, b.Width, b.TitleBarHeight, b.TitleBarColor)
		rl.DrawTextEx(Default_Widget_Header_Font, b.Name, rl.NewVector2(b.TitleBarBounds.X+20, b.TitleBarBounds.Y+5), 14, 0, b.TextColor)
	}
	// Draw border
	rl.DrawRectangleLinesEx(b.TitleBarBounds, b.BorderThickness, b.BorderColor)

	// Draw Window Border
	if b.DrawBorder {
		rl.DrawRectangleLinesEx(b.Bounds, 1, b.BorderColor)
	}
}

func (b *BaseWidget) Update() {

	mouse_pos := rl.GetMousePosition()

	// Check if mouse is over title bar
	isMouseOverTitleBar := rl.CheckCollisionPointRec(mouse_pos, b.TitleBarBounds)

	if isMouseOverTitleBar {
		b.BorderColor = rl.White
		b.BorderThickness = 2

		// Start dragging when left mouse button is pressed
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			b.IsDragging = true
			b.DragOffsetX = mouse_pos.X - b.Bounds.X
			b.DragOffsetY = mouse_pos.Y - b.Bounds.Y
		}
	} else {
		b.BorderColor = rl.Black
		b.BorderThickness = 1
	}

	// Handle dragging
	if b.IsDragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			// Update widget position based on mouse position and offset
			b.Bounds.X = mouse_pos.X - b.DragOffsetX
			b.Bounds.Y = mouse_pos.Y - b.DragOffsetY

			// Keep widget within screen bounds (optional)
			b.keepWithinScreen()
		} else {
			// Stop dragging when mouse button is released
			b.IsDragging = false
		}
	}
	// Update title bar bounds
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		float32(b.Width),
		float32(b.TitleBarHeight),
	)
}

func (b *BaseWidget) keepWithinScreen() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Keep widget within left and right boundaries
	if b.Bounds.X < 0 {
		b.Bounds.X = 0
	} else if b.Bounds.X+float32(b.Width) > screenWidth {
		b.Bounds.X = screenWidth - float32(b.Width)
	}

	// Keep widget within top and bottom boundaries
	if b.Bounds.Y < 0 {
		b.Bounds.Y = 0
	} else if b.Bounds.Y+float32(b.Height) > screenHeight {
		b.Bounds.Y = screenHeight - float32(b.Height)
	}
}
