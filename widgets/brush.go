package widgets

import (
	"go-canvas/globals"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Patterns map[int]rl.Texture2D

const (
	Circle    = 0
	Rectangle = 1
	Round     = 2
	Square    = 3
	Bucket    = 4
)

type Brush struct {
	Shape      int
	Size       float32
	Color      rl.Color
	Pattern    rl.Texture2D
	UsePattern bool
	rotation   float32
}

func NewBrush() *Brush {
	return &Brush{
		Shape:      Circle,
		Size:       5,
		Color:      rl.Black,
		UsePattern: false,
		rotation:   0,
	}
}

func (b *Brush) Update(layer_editor *LayerEditor, color rl.Color) {
	// Only clear strokes once mouse button is released,
	// so that connecting lines work correctly while dragging.
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		for _, layer := range layer_editor.Slots {
			if layer.Visibility {
				// Clear strokes slice without allocating a new one
				layer.Strokes = layer.Strokes[:0]
			}
		}
	}
	b.Color = color
}

func (b *Brush) DrawBrush(layer *LayerSlot, canvas *PaintCanvas) {
	mousepos := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mousepos, canvas.Bounds) {
		return
	}

	if b.Shape == Bucket {
		// Draw a bucket icon or cursor
		rl.DrawRectangleLines(int32(mousepos.X)-5, int32(mousepos.Y)-5, 10, 10, b.Color)
		rl.DrawLine(int32(mousepos.X)-7, int32(mousepos.Y)+7, int32(mousepos.X)+7, int32(mousepos.Y)+7, b.Color)
	} else if b.UsePattern {
		// Draw pattern preview with color tint
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(b.Pattern.Width), Height: float32(b.Pattern.Height)}
		dst := rl.Rectangle{
			X:      mousepos.X,
			Y:      mousepos.Y,
			Width:  float32(b.Size),
			Height: float32(b.Size),
		}
		rl.DrawTexturePro(b.Pattern, src, dst, rl.Vector2{X: float32(b.Size) / 2, Y: float32(b.Size) / 2}, b.rotation, b.Color)
	} else {
		switch b.Shape {
		case Circle:
			rl.DrawCircle(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
		case Rectangle:
			rl.DrawRectangle(int32(mousepos.X), int32(mousepos.Y), int32(b.Size), int32(b.Size), b.Color)
		case Round:
			rl.DrawCircleLines(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
		case Square:
			rl.DrawRectangleLines(int32(mousepos.X-b.Size/2), int32(mousepos.Y-b.Size/2), int32(b.Size), int32(b.Size), b.Color)
		}
	}
}

func (b *Brush) PaintLayer(layer *LayerSlot, canvas *PaintCanvas) {
	if !layer.Visibility {
		return
	}

	mouse := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mouse, canvas.Bounds) {
		return
	}
	localX := mouse.X - canvas.Bounds.X
	localY := mouse.Y - canvas.Bounds.Y

	// Handle paint bucket tool
	if b.Shape == Bucket && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		// Use the flood fill algorithm
		SimpleFloodFill(layer, int(localX), int(localY), b.Color, canvas)
		return
	}

	// Handle other brush tools (existing code)
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) && b.Shape != Bucket {
		// only add if far enough from last point
		const minGap = float32(10) // tune to taste

		if n := len(layer.Strokes); n == 0 ||
			math.Hypot(float64(localX-layer.Strokes[n-1].X), float64(localY-layer.Strokes[n-1].Y)) >= float64(minGap) {

			stroke := globals.Stroke{
				Id:         globals.GenerateUniqueID(),
				X:          localX,
				Y:          localY,
				Size:       float32(b.Size),
				Color:      b.Color,
				Shape:      b.Shape,
				UsePattern: b.UsePattern,
				Pattern:    b.Pattern,
			}
			layer.AddStroke(&stroke)
		}
	}
}

func (b *Brush) DrawStrokes(layer_editor *LayerEditor) {
	for _, layer := range layer_editor.Slots {
		if !layer.Visibility {
			continue
		}
		if len(layer.Strokes) < 1 {
			continue
		}

		rl.BeginTextureMode(layer.Texture)

		for i, stroke := range layer.Strokes {
			x, y := stroke.X, stroke.Y

			if stroke.UsePattern {

				// Pattern strokes with proper alpha handling
				src := rl.Rectangle{X: 0, Y: 0, Width: float32(stroke.Pattern.Width), Height: float32(stroke.Pattern.Height)}
				dst := rl.Rectangle{
					X:      x, // Center the pattern
					Y:      y,
					Width:  stroke.Size,
					Height: stroke.Size,
				}
				rl.DrawTexturePro(
					stroke.Pattern,
					src,
					dst,
					rl.Vector2{X: stroke.Size / 2, Y: stroke.Size / 2}, // Rotate around center
					0,
					stroke.Color,
				)

			} else {
				// Solid shapes: use multiply blending so colors darken each other
				switch stroke.Shape {
				case Circle:
					rl.DrawCircle(int32(x), int32(y), stroke.Size, stroke.Color)
				case Rectangle:
					rl.DrawRectangle(int32(x), int32(y),
						int32(stroke.Size), int32(stroke.Size), stroke.Color)
				case Round:
					rl.DrawCircleLines(int32(x), int32(y),
						float32(stroke.Size), stroke.Color)
				case Square:
					rl.DrawRectangleLines(int32(x-stroke.Size/2), int32(y-stroke.Size/2),
						int32(stroke.Size), int32(stroke.Size), stroke.Color)
				}
			}

			// Connect to previous stroke for smoother lines
			if i > 0 && !stroke.UsePattern && stroke.Shape != Round && stroke.Shape != Square {
				rl.BeginBlendMode(rl.BlendDstRgb)
				prev := layer.Strokes[i-1]
				start := rl.Vector2{X: prev.X, Y: prev.Y}
				end := rl.Vector2{X: x, Y: y}
				if stroke.Shape == Rectangle {
					rl.DrawLineEx(start, end, stroke.Size/5, stroke.Color)
				} else {
					rl.DrawLineEx(start, end, stroke.Size*2, stroke.Color)
				}
				rl.EndBlendMode()
			}
		}

		rl.EndTextureMode()
	}
}
