package widgets

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ColorSet struct {
	BaseWidget
	ColorButton     rl.Rectangle
	ToolBox         *ToolBox
	HightlightColor rl.Color
	CurrentColor    rl.Color

	RRect rl.Rectangle
	GRect rl.Rectangle
	BRect rl.Rectangle
	ARect rl.Rectangle

	R_handle rl.Rectangle
	G_handle rl.Rectangle
	B_handle rl.Rectangle
	A_handle rl.Rectangle

	// RGBA values (0-255)
	RValue int
	GValue int
	BValue int
	AValue int

	// Track which slider is being dragged
	DraggingR bool
	DraggingG bool
	DraggingB bool
	DraggingA bool

	DebugDraw bool
}

func NewColorSet(name string, toolBox *ToolBox) *ColorSet {
	colorset := ColorSet{}
	colorset.ToolBox = toolBox
	// Setup BaseWidget first
	colorset.Visible = true
	colorset.Name = name
	colorset.BodyColor = rl.NewColor(50, 50, 50, 50)
	colorset.DrawTitleBar = false
	colorset.TextColor = rl.White
	colorset.Width = colorset.ToolBox.Width
	colorset.Height = 200

	colorset.Bounds = rl.NewRectangle(
		colorset.ToolBox.Bounds.X,
		colorset.ToolBox.Bounds.Y+colorset.ToolBox.Bounds.Height,
		float32(colorset.Width),
		float32(colorset.Height))

	colorset.HightlightColor = rl.NewColor(255, 255, 255, 255)
	colorset.DebugDraw = false

	// Initialize RGBA values (0-255)
	colorset.RValue = 0   // Red
	colorset.GValue = 255 // Green
	colorset.BValue = 0   // Blue
	colorset.AValue = 255 // Full alpha

	// Set initial color
	colorset.CurrentColor = rl.NewColor(
		uint8(colorset.RValue),
		uint8(colorset.GValue),
		uint8(colorset.BValue),
		uint8(colorset.AValue))

	return &colorset
}

func (c *ColorSet) Draw() {
	c.BaseWidget.Draw()
	c.Update()

	// DrawButton
	rl.DrawRectangle(c.ColorButton.ToInt32().X, c.ColorButton.ToInt32().Y, c.ColorButton.ToInt32().Width, c.ColorButton.ToInt32().Height, c.CurrentColor)

	// Draw R slider
	rl.DrawRectangle(c.RRect.ToInt32().X, c.RRect.ToInt32().Y, c.RRect.ToInt32().Width, c.RRect.ToInt32().Height, rl.LightGray)
	rl.DrawRectangle(c.R_handle.ToInt32().X, c.R_handle.ToInt32().Y, c.R_handle.ToInt32().Width, c.R_handle.ToInt32().Height, rl.Red)
	rl.DrawTextEx(Default_Widget_Header_Font, "R", rl.NewVector2(c.RRect.X-2, c.RRect.Y+c.RRect.Height+5), 14, 0, c.TextColor)
	rl.DrawTextEx(Default_tiny_Font, fmt.Sprintf("%d", c.RValue), rl.NewVector2(c.RRect.X-5, c.RRect.Y-10), 10, 0, c.TextColor)

	// Draw G Slider
	rl.DrawRectangle(c.GRect.ToInt32().X, c.GRect.ToInt32().Y, c.GRect.ToInt32().Width, c.GRect.ToInt32().Height, rl.LightGray)
	rl.DrawRectangle(c.G_handle.ToInt32().X, c.G_handle.ToInt32().Y, c.G_handle.ToInt32().Width, c.G_handle.ToInt32().Height, rl.Green)
	rl.DrawTextEx(Default_Widget_Header_Font, "G", rl.NewVector2(c.GRect.X-2, c.GRect.Y+c.GRect.Height+5), 14, 0, c.TextColor)
	rl.DrawTextEx(Default_tiny_Font, fmt.Sprintf("%d", c.GValue), rl.NewVector2(c.GRect.X-5, c.GRect.Y-10), 10, 0, c.TextColor)

	// Draw B Slider
	rl.DrawRectangle(c.BRect.ToInt32().X, c.BRect.ToInt32().Y, c.BRect.ToInt32().Width, c.BRect.ToInt32().Height, rl.LightGray)
	rl.DrawRectangle(c.B_handle.ToInt32().X, c.B_handle.ToInt32().Y, c.B_handle.ToInt32().Width, c.B_handle.ToInt32().Height, rl.Blue)
	rl.DrawTextEx(Default_Widget_Header_Font, "B", rl.NewVector2(c.BRect.X-2, c.BRect.Y+c.BRect.Height+5), 14, 0, c.TextColor)
	rl.DrawTextEx(Default_tiny_Font, fmt.Sprintf("%d", c.BValue), rl.NewVector2(c.BRect.X-5, c.BRect.Y-10), 10, 0, c.TextColor)

	// Draw A Slider
	rl.DrawRectangle(c.ARect.ToInt32().X, c.ARect.ToInt32().Y, c.ARect.ToInt32().Width, c.ARect.ToInt32().Height, rl.LightGray)
	rl.DrawRectangle(c.A_handle.ToInt32().X, c.A_handle.ToInt32().Y, c.A_handle.ToInt32().Width, c.A_handle.ToInt32().Height, rl.White)
	rl.DrawTextEx(Default_Widget_Header_Font, "A", rl.NewVector2(c.ARect.X-2, c.ARect.Y+c.ARect.Height+5), 14, 0, c.TextColor)
	rl.DrawTextEx(Default_tiny_Font, fmt.Sprintf("%d", c.AValue), rl.NewVector2(c.ARect.X-5, c.ARect.Y-10), 10, 0, c.TextColor)

	// debug draw
	if c.DebugDraw {
		rl.DrawRectangleLinesEx(c.ColorButton, 1, rl.Red)
		rl.DrawRectangleLinesEx(c.RRect, 1, rl.Red)
		rl.DrawRectangleLinesEx(c.GRect, 1, rl.Red)
		rl.DrawRectangleLinesEx(c.BRect, 1, rl.Red)
		rl.DrawRectangleLinesEx(c.ARect, 1, rl.Red)
	}
}

func (c *ColorSet) Update() {
	c.Bounds.X = c.ToolBox.Bounds.X
	c.Bounds.Y = c.ToolBox.Bounds.Y + c.ToolBox.Bounds.Height + 50
	c.Bounds.Width = c.ToolBox.Bounds.Width
	c.ColorButton = rl.NewRectangle(c.Bounds.X+10, c.Bounds.Y+10, c.Bounds.Width-20, c.Bounds.Width-40)

	c.RRect = rl.NewRectangle(c.ColorButton.X, c.ColorButton.Y+c.ColorButton.Height+15, 2, 100)
	c.GRect = rl.NewRectangle(c.ColorButton.X+18, c.ColorButton.Y+c.ColorButton.Height+15, 2, 100)
	c.BRect = rl.NewRectangle(c.ColorButton.X+36, c.ColorButton.Y+c.ColorButton.Height+15, 2, 100)
	c.ARect = rl.NewRectangle(c.ColorButton.X+54, c.ColorButton.Y+c.ColorButton.Height+15, 2, 100)

	// Update handle positions based on current RGBA values
	c.R_handle = rl.NewRectangle(c.RRect.X-5, c.RRect.Y+float32(c.RValue)*c.RRect.Height/255, 12, 10)
	c.G_handle = rl.NewRectangle(c.GRect.X-5, c.GRect.Y+float32(c.GValue)*c.GRect.Height/255, 12, 10)
	c.B_handle = rl.NewRectangle(c.BRect.X-5, c.BRect.Y+float32(c.BValue)*c.BRect.Height/255, 12, 10)
	c.A_handle = rl.NewRectangle(c.ARect.X-5, c.ARect.Y+float32(c.AValue)*c.ARect.Height/255, 12, 10)

	mouse_pos := rl.GetMousePosition()

	// Handle dragging logic
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		if c.DraggingR || rl.CheckCollisionPointRec(mouse_pos, c.R_handle) {
			c.DraggingR = true
			c.updateRValue(mouse_pos)
		} else if c.DraggingG || rl.CheckCollisionPointRec(mouse_pos, c.G_handle) {
			c.DraggingG = true
			c.updateGValue(mouse_pos)
		} else if c.DraggingB || rl.CheckCollisionPointRec(mouse_pos, c.B_handle) {
			c.DraggingB = true
			c.updateBValue(mouse_pos)
		} else if c.DraggingA || rl.CheckCollisionPointRec(mouse_pos, c.A_handle) {
			c.DraggingA = true
			c.updateAValue(mouse_pos)
		}
	} else {
		c.DraggingR = false
		c.DraggingG = false
		c.DraggingB = false
		c.DraggingA = false
	}

	// Update current color based on RGBA values
	c.CurrentColor = rl.NewColor(
		uint8(c.RValue),
		uint8(c.GValue),
		uint8(c.BValue),
		uint8(c.AValue))

	// Draw highlight when hovering over handles
	if rl.CheckCollisionPointRec(mouse_pos, c.R_handle) && !c.DraggingR {
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(
				c.R_handle.X-4, c.R_handle.Y-4, c.R_handle.Width+8, c.R_handle.Height+8), 2, c.HightlightColor)
	}

	if rl.CheckCollisionPointRec(mouse_pos, c.G_handle) && !c.DraggingG {
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(
				c.G_handle.X-4, c.G_handle.Y-4, c.G_handle.Width+8, c.G_handle.Height+8), 2, c.HightlightColor)
	}

	if rl.CheckCollisionPointRec(mouse_pos, c.B_handle) && !c.DraggingB {
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(
				c.B_handle.X-4, c.B_handle.Y-4, c.B_handle.Width+8, c.B_handle.Height+8), 2, c.HightlightColor)
	}

	if rl.CheckCollisionPointRec(mouse_pos, c.A_handle) && !c.DraggingA {
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(
				c.A_handle.X-4, c.A_handle.Y-4, c.A_handle.Width+8, c.A_handle.Height+8), 2, c.HightlightColor)
	}
}

func (c *ColorSet) updateRValue(mouse_pos rl.Vector2) {
	newY := mouse_pos.Y - c.RRect.Y
	if newY < 0 {
		newY = 0
	} else if newY > c.RRect.Height {
		newY = c.RRect.Height
	}

	c.RValue = int(newY * 255 / c.RRect.Height)
	if c.RValue < 0 {
		c.RValue = 0
	} else if c.RValue > 255 {
		c.RValue = 255
	}
}

func (c *ColorSet) updateGValue(mouse_pos rl.Vector2) {
	newY := mouse_pos.Y - c.GRect.Y
	if newY < 0 {
		newY = 0
	} else if newY > c.GRect.Height {
		newY = c.GRect.Height
	}

	c.GValue = int(newY * 255 / c.GRect.Height)
	if c.GValue < 0 {
		c.GValue = 0
	} else if c.GValue > 255 {
		c.GValue = 255
	}
}

func (c *ColorSet) updateBValue(mouse_pos rl.Vector2) {
	newY := mouse_pos.Y - c.BRect.Y
	if newY < 0 {
		newY = 0
	} else if newY > c.BRect.Height {
		newY = c.BRect.Height
	}

	c.BValue = int(newY * 255 / c.BRect.Height)
	if c.BValue < 0 {
		c.BValue = 0
	} else if c.BValue > 255 {
		c.BValue = 255
	}
}

func (c *ColorSet) updateAValue(mouse_pos rl.Vector2) {
	newY := mouse_pos.Y - c.ARect.Y
	if newY < 0 {
		newY = 0
	} else if newY > c.ARect.Height {
		newY = c.ARect.Height
	}

	c.AValue = int(newY * 255 / c.ARect.Height)
	if c.AValue < 0 {
		c.AValue = 0
	} else if c.AValue > 255 {
		c.AValue = 255
	}
}

func (c *ColorSet) OnColorClick() {
	// This function can be used for additional color click functionality
}
