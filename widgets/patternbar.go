package widgets

import (
	"go-canvas/globals"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BrushPatternBar struct {
	BaseWidget
	Canvas         *PaintCanvas
	Patterns       []rl.Rectangle
	CurrentPattern rl.Rectangle
	CurrentTexture *rl.Texture2D

	BrushColor      rl.Color
	HightlightColor rl.Color
	DebugDraw       bool
}

func NewBrushPatternBar(name string, canvas *PaintCanvas) *BrushPatternBar {
	pattern_bar := BrushPatternBar{}
	pattern_bar.Canvas = canvas
	// Setup BaseWidget first
	pattern_bar.BrushColor = rl.White
	pattern_bar.Visible = true
	pattern_bar.Name = name
	pattern_bar.DrawTitleBar = false
	pattern_bar.TextColor = rl.White
	pattern_bar.BodyColor = rl.NewColor(50, 50, 50, 20)
	pattern_bar.Width = 680
	pattern_bar.Height = 70
	pattern_bar.BorderThickness = 1
	pattern_bar.BorderColor = rl.Green
	pattern_bar.Bounds = rl.NewRectangle(10, float32(rl.GetScreenHeight()-200), float32(pattern_bar.Width), pattern_bar.Bounds.Height)
	pattern_bar.HightlightColor = rl.NewColor(255, 255, 255, 100)
	pattern_bar.DebugDraw = false

	return &pattern_bar
}

func (p *BrushPatternBar) Draw() {
	p.BaseWidget.Draw()
	p.Update()

	i := float32(0)
	p.Patterns = nil
	for _, texture := range globals.BRUSH_PATTERNS {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}
		dest := rl.NewRectangle(p.Bounds.X+i, p.Bounds.Y, 50, 50)
		p.Patterns = append(p.Patterns, dest)
		rl.DrawTexturePro(
			texture,
			src,
			dest,
			rl.Vector2{X: -10, Y: -10},
			0,
			p.BrushColor)
		i += 55
	}
	rl.DrawRectangleLinesEx(rl.NewRectangle(p.CurrentPattern.X+7, p.CurrentPattern.Y+10, p.CurrentPattern.Width, p.CurrentPattern.Height), 1, rl.White)

}

func (p *BrushPatternBar) Update() {
	p.Bounds = rl.NewRectangle(float32(rl.GetScreenWidth()/2)-310, float32(rl.GetScreenHeight()-100), 500, 50)

	mouse_pos := rl.GetMousePosition()
	for i := 0; i < len(p.Patterns); i++ {
		if rl.CheckCollisionPointRec(mouse_pos, p.Patterns[i]) {
			rl.DrawRectangle(int32(p.Patterns[i].X+7), int32(p.Patterns[i].Y+10), int32(p.Patterns[i].Width), int32(p.Patterns[i].Height), p.HightlightColor)
		}
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mouse_pos, p.Patterns[i]) {
			p.CurrentPattern = p.Patterns[i]
			p.CurrentTexture = &globals.BRUSH_PATTERNS[i]
			p.Canvas.Brush.Pattern = globals.BRUSH_PATTERNS[i]
			p.Canvas.Brush.UsePattern = true
		}
	}

}
