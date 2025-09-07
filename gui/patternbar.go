package widgets

import (
	"fmt"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BrushPatternBar struct {
	BaseWidget
	Patterns       []rl.Rectangle
	BrushPatterns  []rl.Texture2D
	CurrentPattern rl.Rectangle
	CurrentTexture *rl.Texture2D

	BrushColor      rl.Color
	HightlightColor rl.Color
	DebugDraw       bool
}

func NewBrushPatternBar(name string) *BrushPatternBar {
	pattern_bar := BrushPatternBar{}
	pattern_bar.BrushPatterns = make([]rl.Texture2D, 0)
	pattern_bar.LoadToolBoxIcons()
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
	for _, texture := range p.BrushPatterns {
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
			p.CurrentTexture = &p.BrushPatterns[i]
		}
	}

}

func (p *BrushPatternBar) LoadToolBoxIcons() {
	entries, err := os.ReadDir("./patterns")
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for i < 12 {
		texture := rl.LoadTexture(fmt.Sprintf("patterns/%v", entries[i].Name()))
		p.BrushPatterns = append(p.BrushPatterns, texture)
		i += 1
		if i >= len(entries) {
			break
		}

	}

}
