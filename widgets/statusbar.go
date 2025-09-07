package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type StatusBar struct {
	BaseWidget
}

func NewStatusBar(name string) *StatusBar {
	status_bar := StatusBar{}

	// Setup BaseWidget first
	status_bar.Visible = true
	status_bar.Name = name
	status_bar.BodyColor = rl.NewColor(50, 50, 50, 30)
	status_bar.DrawTitleBar = false
	status_bar.TextColor = rl.White
	status_bar.Width = int32(rl.GetScreenWidth()) - 15
	status_bar.Height = 15

	status_bar.Bounds = rl.NewRectangle(10, float32(rl.GetScreenHeight())-25, float32(status_bar.Width), float32(status_bar.Height))

	return &status_bar
}

func (s *StatusBar) Draw() {
	s.BaseWidget.Draw()
	s.Update()

	// draw shot-cut text

	text := "[SHORT-CUTS]  'B' : Brush Mode,  'P' : Pattern Mode,  'E' : Eraser,  'ctl+S' : Save,  'ctl+N' : New,  'ctl+R' : clear"
	rl.DrawTextEx(Default_Widget_Header_Font, text, rl.NewVector2(s.Bounds.X, s.Bounds.Y), 14, 0, rl.White)
}
func (s *StatusBar) Update() {
	s.Bounds.Y = float32(rl.GetScreenHeight() - 25)
	s.Width = int32(rl.GetScreenWidth()) - 15

}
