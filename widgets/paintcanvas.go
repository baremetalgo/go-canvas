package widgets

import (
	"go-canvas/globals"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PaintCanvas struct {
	BaseWidget
	Texture         rl.RenderTexture2D
	LayerEditor     *LayerEditor
	Brush           *Brush
	Active_Layer_Id float32
}

func NewPaintCanvas(name string, layer_editor *LayerEditor) *PaintCanvas {
	canvas := PaintCanvas{}
	canvas.LayerEditor = layer_editor
	canvas.Brush = NewBrush()
	canvas.Texture = rl.LoadRenderTexture(
		globals.CANVAS_WIDTH, globals.CANVAS_HEIGHT)
	canvas.Visible = true
	canvas.Name = name
	canvas.BodyColor = rl.NewColor(50, 50, 50, 30)
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

	// rl.DrawRectanglePro(s.Bounds, rl.NewVector2(0, 0), 0, rl.White)

}
func (s *PaintCanvas) Update() {
	mouse_pos := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mouse_pos, s.Bounds) {
		rl.SetMouseCursor(rl.MouseCursorPointingHand)
	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}

}

func (c *PaintCanvas) CompositeLayers() {
	// Composite all layers into the final texture
	rl.BeginTextureMode(rl.RenderTexture2D(c.Texture))
	// rl.ClearBackground(rl.Blank)

	for _, layer := range c.LayerEditor.Slots {
		if !layer.Visibility {
			continue
		}

		rl.BeginBlendMode(BlendModes[layer.BlendMode])
		// Draw the layer texture with its opacity
		rl.DrawTextureRec(
			layer.Texture.Texture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(c.Width), Height: -float32(c.Height)},
			rl.Vector2{X: 0, Y: 0},
			rl.Fade(rl.White, layer.Opacity),
		)
		rl.EndBlendMode()
	}

	rl.EndTextureMode()
}
