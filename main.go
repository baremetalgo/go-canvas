package main

import (
	"go-canvas/globals"
	"go-canvas/gui"
	"go-canvas/widgets"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboardInputs(brush *widgets.Brush) {
	if rl.IsKeyDown(rl.KeyKpAdd) {
		brush.Size += 0.1
	}

	if rl.IsKeyDown(rl.KeyKpSubtract) {
		if brush.Size > 1 {
			brush.Size -= 0.1
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		if brush.Shape == widgets.Circle {
			brush.Shape = widgets.Rectangle
		} else {
			brush.Shape = widgets.Circle
		}
	}
	if rl.IsKeyPressed(rl.KeyS) {
		rl.TakeScreenshot("painting.png")
	}
	if rl.IsKeyPressed(rl.KeyP) {
		brush.UsePattern = !brush.UsePattern
		if brush.UsePattern {
			brush.Pattern = widgets.Patterns[0] // pick first pattern
		}
	}
}

func init_raylib_window() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1024, 786, "Go-Painter!")

}
func main() {
	init_raylib_window()
	defer rl.CloseWindow()
	globals.LoadToolBoxIcons()

	// widgets.LoadPatterns()
	// defer func() {
	// 	for _, tex := range widgets.Patterns {
	// 		rl.UnloadTexture(tex)
	// 	}
	// }()

	widgets.InitializeFonts()

	ui := gui.NewUserInterface()
	canvas_texture := ui.PaintCanvas.Texture.Texture

	for !rl.WindowShouldClose() {
		active_layer := ui.PaintCanvas.LayerEditor.ActiveLayer
		HandleKeyboardInputs(ui.PaintCanvas.Brush)

		// Paint new strokes to the layer texture
		ui.PaintCanvas.Brush.PaintLayer(active_layer, ui.PaintCanvas)
		ui.PaintCanvas.Brush.DrawStrokes(ui.PaintCanvas.LayerEditor)

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		ui.Draw()
		// Draw the layer texture
		ui.PaintCanvas.CompositeLayers()
		tex_x := float32(rl.GetScreenWidth())/2 - float32(canvas_texture.Width)/2
		tex_y := float32(rl.GetScreenHeight())/2 - float32(canvas_texture.Height)/2

		src := rl.NewRectangle(0, 0, float32(canvas_texture.Width), float32(canvas_texture.Height))
		dst := rl.NewRectangle(tex_x, tex_y, float32(canvas_texture.Width), float32(canvas_texture.Height))

		rl.DrawTexturePro(
			canvas_texture,
			src,
			dst,
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)

		// Draw brush preview
		ui.PaintCanvas.Brush.DrawBrush(active_layer)

		rl.DrawFPS(100, 100)
		rl.EndDrawing()
	}
}
