package main

import (
	"go-canvas/gui"
	"go-canvas/widgets"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboardInputs(brush *widgets.Brush) {
	if rl.IsKeyDown(rl.KeyKpAdd) {
		brush.Size += 1
	}

	if rl.IsKeyDown(rl.KeyKpSubtract) {
		if brush.Size > 1 {
			brush.Size -= 1
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

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 600, "Go-Painter!")
	defer rl.CloseWindow()

	widgets.LoadPatterns()
	defer func() {
		for _, tex := range widgets.Patterns {
			rl.UnloadTexture(tex)
		}
	}()
	widgets.InitializeFonts()

	ui := gui.NewUserInterface()

	for !rl.WindowShouldClose() {
		active_layer := ui.PaintCanvas.LayerEditor.ActiveLayer
		HandleKeyboardInputs(ui.PaintCanvas.Brush)

		// Paint new strokes to the layer texture
		ui.PaintCanvas.Brush.PaintLayer(active_layer)
		ui.PaintCanvas.Brush.DrawStrokes(ui.PaintCanvas.LayerEditor)

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)

		// Draw the layer texture
		ui.PaintCanvas.CompositeLayers()
		text := ui.PaintCanvas.Texture.Texture
		src := rl.NewRectangle(0, 0, float32(text.Width), float32(text.Height))
		dst := rl.NewRectangle(100, 100, float32(text.Width), float32(text.Height))
		rl.DrawTexturePro(
			text,
			src,
			dst,
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.Red,
		)
		ui.Draw()
		// Draw brush preview
		ui.PaintCanvas.Brush.DrawBrush(active_layer)

		rl.DrawFPS(100, 100)
		rl.EndDrawing()
	}
}
