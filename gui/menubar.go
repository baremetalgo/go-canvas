package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var MenuIcons map[string]rl.Texture2D

type MenuBar struct {
	BaseWidget
	NewButton  rl.Rectangle
	OpenButton rl.Rectangle
	SaveButton rl.Rectangle

	HightlightColor rl.Color
	DebugDraw       bool

	button_map   map[rl.Texture2D]*rl.Rectangle
	button_order [3]*rl.Rectangle
}

func NewMenuBar(name string) *MenuBar {
	menubar := MenuBar{}
	menubar.LoadMenuIcons()

	// Setup BaseWidget first
	menubar.Visible = true
	menubar.Name = name
	menubar.BodyColor = rl.Gray
	menubar.DrawTitleBar = false
	menubar.TextColor = rl.White
	menubar.Width = 200
	menubar.Height = 100
	menubar.BorderThickness = 1
	menubar.BorderColor = rl.Green
	menubar.Bounds = rl.NewRectangle(float32(rl.GetScreenWidth()/2)-100, 10, float32(menubar.Width), menubar.Bounds.Height)
	menubar.HightlightColor = rl.NewColor(255, 255, 255, 200)
	menubar.DebugDraw = false
	// Create brush map
	menubar.button_map = map[rl.Texture2D]*rl.Rectangle{
		MenuIcons["new"]:  &menubar.NewButton,
		MenuIcons["open"]: &menubar.OpenButton,
		MenuIcons["save"]: &menubar.SaveButton,
	}
	menubar.button_order = [3]*rl.Rectangle{&menubar.NewButton, &menubar.OpenButton, &menubar.SaveButton}
	return &menubar
}

func (m *MenuBar) Draw() {
	// m.BaseWidget.Draw()
	m.Update()

	for text, button := range m.button_map {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(text.Width), Height: float32(text.Height)}
		dest := rl.NewRectangle(button.X, button.Y, button.Width, button.Height)
		rl.DrawTexturePro(
			text,
			src,
			dest,
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White)

	}
	// debug draw
	if m.DebugDraw {
		rl.DrawRectangleLinesEx(m.NewButton, 1, rl.Red)
		rl.DrawRectangleLinesEx(m.OpenButton, 1, rl.Red)
		rl.DrawRectangleLinesEx(m.SaveButton, 1, rl.Red)
		rl.DrawRectangleLinesEx(m.Bounds, 1, rl.Red)
	}

}

func (m *MenuBar) Update() {
	m.Bounds = rl.NewRectangle(float32(rl.GetScreenWidth()/2)-100, 10, float32(m.Width), m.Bounds.Height)

	prev_button_x := float32(0)
	overall_width := int32(0)

	mouse_pos := rl.GetMousePosition()

	for _, button := range m.button_order {
		button.X = m.Bounds.X + 10 + prev_button_x
		button.Y = m.Bounds.Y
		button.Width = 30
		button.Height = 30

		prev_button_x += button.Width
		overall_width += button.ToInt32().Width
		if rl.CheckCollisionPointRec(mouse_pos, *button) {
			rl.DrawRectangleLinesEx(*button, 2, m.HightlightColor)
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				rl.DrawRectangleLinesEx(*button, 20, rl.Green)
			}
		}
	}

	m.Width = int32(overall_width)

}

func (m *MenuBar) LoadMenuIcons() {
	MenuIcons = map[string]rl.Texture2D{
		"new":  rl.LoadTexture("icons/new_icon.png"),
		"open": rl.LoadTexture("icons/open_icon.png"),
		"save": rl.LoadTexture("icons/save_icon.png"),
	}
}
