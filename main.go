package main

import (
	"mspaint_golang/core"
	widgets "mspaint_golang/gui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboardInputs(brush *core.Brush, layer *core.Layer) {
	if rl.IsKeyDown(rl.KeyKpAdd) {
		brush.Size += 1
	}

	if rl.IsKeyDown(rl.KeyKpSubtract) {
		if brush.Size > 1 {
			brush.Size -= 1
		}
	}
	if rl.IsKeyPressed(rl.KeyX) {
		layer.Strokes = nil
		// Clear the texture
		rl.BeginTextureMode(layer.Texture)
		rl.ClearBackground(rl.Blank)
		rl.EndTextureMode()
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		if brush.Shape == core.Circle {
			brush.Shape = core.Rectangle
		} else {
			brush.Shape = core.Circle
		}
	}
	if rl.IsKeyPressed(rl.KeyS) {
		rl.TakeScreenshot("painting.png")
	}
	if rl.IsKeyPressed(rl.KeyP) {
		brush.UsePattern = !brush.UsePattern
		if brush.UsePattern {
			brush.Pattern = core.Patterns[0] // pick first pattern
		}
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 600, "Go-Painter!")
	defer rl.CloseWindow()

	core.LoadPatterns()
	defer func() {
		for _, tex := range core.Patterns {
			rl.UnloadTexture(tex)
		}
	}()
	widgets.InitializeFonts()
	canvas := core.NewCanvas()
	widgets := widgets.NewWidgets()

	for !rl.WindowShouldClose() {
		active_layer := canvas.GetActiveLayer()
		HandleKeyboardInputs(canvas.Brush, active_layer)

		// Paint new strokes to the layer texture
		canvas.Brush.PaintLayer(active_layer)
		canvas.Brush.DrawStrokes(active_layer)

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)

		// Draw the layer texture
		rl.DrawTextureRec(
			active_layer.Texture.Texture,
			rl.Rectangle{
				X: 0, Y: 0,
				Width:  float32(active_layer.Texture.Texture.Width),
				Height: -float32(active_layer.Texture.Texture.Height),
			},
			rl.Vector2{X: 0, Y: 0},
			rl.White,
		)

		// Draw brush preview
		// canvas.Brush.DrawBrush()
		widgets.Draw()
		rl.DrawFPS(100, 100)
		rl.EndDrawing()
	}
}
