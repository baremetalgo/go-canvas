package widgets

import (
	"go-canvas/globals"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Patterns map[int]rl.Texture2D

const (
	Circle    = 0
	Rectangle = 1
	Round     = 2
	Square    = 3
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

func (b *Brush) DrawBrush(layer *LayerSlot) {

	mousepos := rl.GetMousePosition()
	if b.UsePattern {
		// Draw pattern preview
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(b.Pattern.Width), Height: float32(b.Pattern.Height)}
		dst := rl.Rectangle{
			X:      mousepos.X - float32(b.Size)/2,
			Y:      mousepos.Y - float32(b.Size)/2,
			Width:  float32(b.Size),
			Height: float32(b.Size),
		}
		rl.DrawTexturePro(b.Pattern, src, dst, rl.Vector2{X: float32(b.Size) / 2, Y: float32(b.Size) / 2}, b.rotation, rl.White)
	} else if b.Shape == Circle {
		rl.DrawCircle(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
	} else if b.Shape == Rectangle {
		rl.DrawRectangle(int32(mousepos.X), int32(mousepos.Y), int32(b.Size), int32(b.Size), b.Color)
	} else if b.Shape == Round {
		rl.DrawCircleLines(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
	} else if b.Shape == Square {
		rl.DrawRectangleLines(int32(mousepos.X), int32(mousepos.Y), int32(b.Size), int32(b.Size), b.Color)
	}
}

func (b *Brush) PaintLayer(layer *LayerSlot, canvas *PaintCanvas) {
	if !layer.Visibility {
		return
	}

	mouse := rl.GetMousePosition()
	localX := mouse.X - canvas.Bounds.X
	localY := mouse.Y - canvas.Bounds.Y

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		// only add if far enough from last point
		const minGap = float32(1) // tune to taste
		if n := len(layer.Strokes); n == 0 ||
			math.Hypot(float64(localX-layer.Strokes[n-1].X), float64(localY-layer.Strokes[n-1].Y)) >= float64(minGap) {

			b.rotation = float32(rand.Intn(361))
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

		// fmt.Printf("Drawing Layer %v, visibility %v, Strokes : %v\n",
		// 	layer.Name, layer.Visibility, len(layer.Strokes))

		rl.BeginTextureMode(layer.Texture)

		for i, stroke := range layer.Strokes {
			x, y := stroke.X, stroke.Y

			if stroke.UsePattern {
				// Draw texture scaled to brush size
				src := rl.Rectangle{X: 0, Y: 0, Width: float32(stroke.Pattern.Width), Height: float32(stroke.Pattern.Height)}
				dst := rl.Rectangle{
					X:      x - stroke.Size/2,
					Y:      y - stroke.Size/2,
					Width:  stroke.Size,
					Height: stroke.Size,
				}
				rl.DrawTexturePro(stroke.Pattern, src, dst,
					rl.Vector2{X: stroke.Size / 2, Y: stroke.Size / 2}, 0, stroke.Color)

			} else {
				// Solid shapes
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
					rl.DrawRectangleLines(int32(x), int32(y),
						int32(stroke.Size), int32(stroke.Size), stroke.Color)
				}

			}

			// Draw connecting line to previous stroke
			if i > 0 {
				prev := layer.Strokes[i-1]
				dx := x - prev.X
				dy := y - prev.Y
				dist := math.Sqrt(float64(dx*dx + dy*dy))

				if stroke.UsePattern {
					// For patterns, draw intermediate points
					step := float32(math.Max(1, float64(stroke.Size)/2))
					steps := int(float32(dist)/step) + 1
					for s := 1; s <= steps; s++ {
						t := float32(s) / float32(steps)
						ix := prev.X + dx*t
						iy := prev.Y + dy*t

						src := rl.Rectangle{X: 0, Y: 0, Width: float32(stroke.Pattern.Width), Height: float32(stroke.Pattern.Height)}
						dst := rl.Rectangle{
							X:      ix - stroke.Size/2,
							Y:      iy - stroke.Size/2,
							Width:  stroke.Size,
							Height: stroke.Size,
						}
						if stroke.Shape == Rectangle {
							rl.DrawTexturePro(stroke.Pattern, src, dst,
								rl.Vector2{X: stroke.Size, Y: stroke.Size / 2}, 0, stroke.Color)
						}
						if stroke.Shape == Circle {
							rl.DrawTexturePro(stroke.Pattern, src, dst,
								rl.Vector2{X: stroke.Size / 2, Y: stroke.Size / 2}, 0, stroke.Color)
						}

					}
				} else {
					// For solid shapes, draw a line
					start := rl.Vector2{X: prev.X, Y: prev.Y}
					end := rl.Vector2{X: x, Y: y}
					rl.DrawLineEx(start, end, stroke.Size*2, stroke.Color)
				}
			}
		}

		rl.EndTextureMode()
	}
}
