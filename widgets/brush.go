package widgets

import (
	"fmt"
	"go-canvas/globals"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Patterns map[int]rl.Texture2D

const (
	Circle    = 0
	Rectangle = 1
)

type Brush struct {
	Shape      int
	Size       int32
	Color      rl.Color
	Pattern    rl.Texture2D
	UsePattern bool
	rotation   float32
}

func NewBrush() *Brush {
	return &Brush{
		Shape:      Circle,
		Size:       5,
		Color:      rl.Red,
		UsePattern: false,
		rotation:   0,
	}
}

func (b *Brush) DrawBrush(layer *LayerSlot) {

	mousepos := rl.GetMousePosition()
	if b.UsePattern {
		// Draw pattern preview
		src := rl.Rectangle{0, 0, float32(b.Pattern.Width), float32(b.Pattern.Height)}
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
		rl.DrawRectangle(int32(mousepos.X), int32(mousepos.Y), b.Size, b.Size, b.Color)
	}
}

func (b *Brush) PaintLayer(layer *LayerSlot) {
	if !layer.Visibility {
		return
	}
	mousepos := rl.GetMousePosition()

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		// Save stroke info including pattern
		b.rotation = float32(rand.Intn(361))
		stroke := globals.Stroke{
			Id:         globals.GenerateUniqueID(),
			X:          mousepos.X,
			Y:          mousepos.Y,
			Size:       float32(b.Size),
			Color:      b.Color,
			Shape:      b.Shape,
			UsePattern: b.UsePattern,
			Pattern:    b.Pattern,
		}
		layer.AddStroke(&stroke)
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

		fmt.Printf("Drawing Layer %v, visibility %v, Strokes : %v\n", layer.Name, layer.Visibility, len(layer.Strokes))
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
				rl.DrawTexturePro(stroke.Pattern, src, dst, rl.Vector2{X: stroke.Size / 2, Y: stroke.Size / 2}, 0, rl.White)
			} else {
				// Solid shapes
				switch stroke.Shape {
				case Circle:
					rl.DrawCircle(int32(x), int32(y), stroke.Size, stroke.Color)
				case Rectangle:
					rl.DrawRectangle(int32(x), int32(y), int32(stroke.Size), int32(stroke.Size), stroke.Color)
				}
			}
			// Draw connecting line if distance < 20
			if i > 0 {
				prev := layer.Strokes[i-1]
				dx := x - prev.X
				dy := y - prev.Y
				dist := math.Sqrt(float64(dx*dx + dy*dy))

				if dist < 20 {
					if stroke.UsePattern {
						// For patterns, draw intermediate points
						steps := int(dist / 2)
						if steps > 0 {
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
								rl.DrawTexturePro(stroke.Pattern, src, dst, rl.Vector2{X: stroke.Size / 2, Y: stroke.Size / 2}, 0, rl.White)
							}
						}
					} else {
						// For solid shapes, draw a line
						start := rl.Vector2{X: prev.X, Y: prev.Y}
						end := rl.Vector2{X: x, Y: y}
						rl.DrawLineEx(start, end, stroke.Size, stroke.Color)
					}
				}
			}
		}
		rl.EndTextureMode()
	}

	for _, layer := range layer_editor.Slots {
		if layer.Visibility {
			layer.Strokes = make([]*globals.Stroke, 0)
		}
	}
}

func LoadPatterns() {
	Patterns = map[int]rl.Texture2D{
		0: rl.LoadTexture("patterns/pencil_pattern.png"),
		1: rl.LoadTexture("sources/W_pressed.png"),
		2: rl.LoadTexture("sources/A_pressed.png"),
		3: rl.LoadTexture("sources/S_pressed.png"),
	}
}
