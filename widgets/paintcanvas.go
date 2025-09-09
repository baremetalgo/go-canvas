package widgets

import (
	"fmt"
	"go-canvas/globals"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PaintCanvas struct {
	BaseWidget
	Texture         rl.RenderTexture2D
	LayerEditor     *LayerEditor
	Brush           *Brush
	ColorSet        *ColorSet
	Active_Layer_Id float32
}

func NewPaintCanvas(name string, layer_editor *LayerEditor, color_set *ColorSet) *PaintCanvas {
	canvas := PaintCanvas{}
	canvas.LayerEditor = layer_editor
	canvas.Brush = NewBrush()
	canvas.ColorSet = color_set
	canvas.Texture = rl.LoadRenderTexture(
		globals.CANVAS_WIDTH, globals.CANVAS_HEIGHT)
	canvas.Texture.Texture.Format = rl.PixelFormat(8)
	canvas.Visible = true
	canvas.Name = name
	canvas.BodyColor = rl.White
	canvas.DrawTitleBar = false
	canvas.TextColor = rl.White
	canvas.Width = globals.CANVAS_WIDTH
	canvas.Height = globals.CANVAS_HEIGHT

	canvas.Bounds = rl.NewRectangle(100, 100, float32(canvas.Width), float32(canvas.Height))

	return &canvas
}

func (s *PaintCanvas) Draw() {
	s.BaseWidget.Draw()
	s.Update()

	rl.DrawTextEx(Default_Widget_Header_Font, fmt.Sprintf("%v x %v", s.Width, s.Height), rl.NewVector2(s.Bounds.X+float32(s.Width)-100, s.Bounds.Y-15), 14, 1, rl.Black)
}

func (s *PaintCanvas) Update() {
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())
	// Center the rectangle
	s.Bounds = rl.NewRectangle(
		float32(screenWidth)/2-float32(s.Width)/2,
		float32(screenHeight)/2-float32(s.Height)/2,
		float32(s.Width),
		float32(s.Height),
	)

	mouse_pos := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mouse_pos, s.Bounds) {
		if s.Brush.Shape != Line {
			rl.SetMouseCursor(rl.MouseCursorPointingHand)
		}
	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}
	s.Brush.Update(s.LayerEditor, s.ColorSet.CurrentColor)

}

func (c *PaintCanvas) CompositeLayers() {
	// Composite all layers into the final texture
	rl.BeginTextureMode(rl.RenderTexture2D(c.Texture))
	rl.ClearBackground(rl.White)

	// Iterate through layers in reverse order (from last to first)
	for i := len(c.LayerEditor.Slots) - 1; i >= 0; i-- {
		layer := c.LayerEditor.Slots[i]
		if !layer.Visibility {
			continue
		}
		if layer.BlendToggle {
			rl.BeginBlendMode(BlendModes[layer.BlendMode])
			// Draw the layer texture with its opacity
			rl.DrawTextureRec(
				layer.Texture.Texture,
				rl.Rectangle{X: 0, Y: 0, Width: float32(c.Width), Height: float32(c.Height)},
				rl.Vector2{X: 0, Y: 0},
				rl.Fade(rl.White, layer.Opacity),
			)
			rl.EndBlendMode()
		} else {
			rl.DrawTextureRec(
				layer.Texture.Texture,
				rl.Rectangle{X: 0, Y: 0, Width: float32(c.Width), Height: float32(c.Height)},
				rl.Vector2{X: 0, Y: 0},
				rl.Fade(rl.White, layer.Opacity),
			)
		}
	}

	rl.EndTextureMode()
}
