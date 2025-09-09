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
	Line      = 5
)

type Brush struct {
	Shape         int
	Size          float32
	Color         rl.Color
	Pattern       rl.Texture2D
	UsePattern    bool
	rotation      float32
	LineStart     rl.Vector2 // Store line start position
	IsDrawingLine bool       // Track if we're in the middle of drawing a line
	TempLineEnd   rl.Vector2 // Temporary end point for preview
}

func NewBrush() *Brush {
	return &Brush{
		Shape:         Circle,
		Size:          5,
		Color:         rl.Black,
		UsePattern:    false,
		rotation:      0,
		IsDrawingLine: false,
	}
}

func (b *Brush) Update(layer_editor *LayerEditor, color rl.Color) {
	// Only clear strokes once mouse button is released,
	// so that connecting lines work correctly while dragging.
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) && b.Shape != Line {
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

	if b.UsePattern {
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
		case Bucket:
			rl.DrawRectangleLines(int32(mousepos.X)-5, int32(mousepos.Y)-5, 30, 30, b.Color)
			rl.DrawLine(int32(mousepos.X)-7, int32(mousepos.Y)+7, int32(mousepos.X)+20, int32(mousepos.Y)+20, b.Color)
		case Circle:
			rl.DrawCircle(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
		case Rectangle:
			rl.DrawRectangle(int32(mousepos.X-b.Size/2), int32(mousepos.Y-b.Size/2), int32(b.Size), int32(b.Size), b.Color)
		case Round:
			rl.DrawCircleLines(int32(mousepos.X), int32(mousepos.Y), float32(b.Size), b.Color)
		case Square:
			rl.DrawRectangleLines(int32(mousepos.X-b.Size/2), int32(mousepos.Y-b.Size/2), int32(b.Size), int32(b.Size), b.Color)
		case Line:
			rl.SetMouseCursor(rl.MouseCursorCrosshair)

			// Draw line preview if we're in the middle of drawing a line
			if b.IsDrawingLine {
				// Convert canvas coordinates to screen coordinates
				screenStart := rl.Vector2{
					X: b.LineStart.X + canvas.Bounds.X,
					Y: b.LineStart.Y + canvas.Bounds.Y,
				}
				rl.DrawLineEx(screenStart, mousepos, b.Size, b.Color)

				// Draw circles at start and end points for better visibility
				rl.DrawCircle(int32(screenStart.X), int32(screenStart.Y), b.Size/2, b.Color)
				rl.DrawCircle(int32(mousepos.X), int32(mousepos.Y), b.Size/2, b.Color)
			} else {
				// Draw crosshair for line tool
				rl.DrawLine(int32(mousepos.X-10), int32(mousepos.Y), int32(mousepos.X+10), int32(mousepos.Y), b.Color)
				rl.DrawLine(int32(mousepos.X), int32(mousepos.Y-10), int32(mousepos.X), int32(mousepos.Y+10), b.Color)
			}
		}
	}
}

func (b *Brush) PaintLayer(layer *LayerSlot, canvas *PaintCanvas) {
	if !layer.Visibility {
		return
	}

	mouse := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mouse, canvas.Bounds) {
		// Cancel line drawing if mouse leaves canvas
		if b.Shape == Line && b.IsDrawingLine {
			b.IsDrawingLine = false
		}
		return
	}
	localX := mouse.X - canvas.Bounds.X
	localY := mouse.Y - canvas.Bounds.Y

	// Handle paint bucket tool
	if b.Shape == Bucket && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && !b.UsePattern {
		// Use the flood fill algorithm
		SimpleFloodFill(layer, int(localX), int(localY), b.Color, canvas)
		return
	}

	// Handle line tool
	if b.Shape == Line {
		b.handleLineTool(layer, canvas, localX, localY)
		return
	}

	// Handle other brush tools (existing code)
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) && b.Shape != Bucket {
		// only add if far enough from last point
		const minGap = float32(10)

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

func (b *Brush) handleLineTool(layer *LayerSlot, canvas *PaintCanvas, localX, localY float32) {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		if !b.IsDrawingLine {
			// Start drawing a new line
			b.IsDrawingLine = true
			b.LineStart = rl.Vector2{X: localX, Y: localY}
			b.TempLineEnd = rl.Vector2{X: localX, Y: localY}
		} else {
			// Finish drawing the line - add it to the layer
			b.IsDrawingLine = false
			b.addLineToLayer(layer, localX, localY)
		}
	} else if b.IsDrawingLine && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		// Update the temporary end point while dragging
		b.TempLineEnd = rl.Vector2{X: localX, Y: localY}
	}

	// Cancel line drawing with right click or escape key
	if (rl.IsMouseButtonPressed(rl.MouseButtonRight) || rl.IsKeyPressed(rl.KeyEscape)) && b.IsDrawingLine {
		b.IsDrawingLine = false
	}
}

func (b *Brush) addLineToLayer(layer *LayerSlot, endX, endY float32) {
	// Create a special stroke that represents a complete line
	// You'll need to modify your Stroke struct to support lines properly
	/*
		stroke := globals.Stroke{
			Id:         globals.GenerateUniqueID(),
			X:          b.LineStart.X,
			Y:          b.LineStart.Y,
			Size:       float32(b.Size),
			Color:      b.Color,
			Shape:      Line,
			UsePattern: b.UsePattern,
			Pattern:    b.Pattern,
			// Add these fields to your Stroke struct in globals package:
			IsLine: true,
			EndX:   endX,
			EndY:   endY,
		}
	*/
	// For now, we'll use a simple approach: store both points as separate strokes
	// but mark them as line endpoints
	startStroke := globals.Stroke{
		Id:    globals.GenerateUniqueID(),
		X:     b.LineStart.X,
		Y:     b.LineStart.Y,
		Size:  float32(b.Size),
		Color: b.Color,
		Shape: Line,
		// IsLine: true,
	}

	endStroke := globals.Stroke{
		Id:    globals.GenerateUniqueID(),
		X:     endX,
		Y:     endY,
		Size:  float32(b.Size),
		Color: b.Color,
		Shape: Line,
		// IsLine: true,
	}

	layer.AddStroke(&startStroke)
	layer.AddStroke(&endStroke)
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

		// First pass: draw all non-line strokes
		for i, stroke := range layer.Strokes {
			if stroke.Shape == Line {
				continue // Skip line strokes for now, we'll handle them separately
			}

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
				// Solid shapes
				switch stroke.Shape {
				case Circle:
					rl.DrawCircle(int32(x), int32(y), stroke.Size, stroke.Color)
				case Rectangle:
					rl.DrawRectangle(int32(x)-int32(stroke.Size)/2, int32(y)-int32(stroke.Size)/2,
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
			if i > 0 && !stroke.UsePattern && stroke.Shape != Round && stroke.Shape != Square && stroke.Shape != Line {
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

		// Second pass: draw lines by connecting consecutive line strokes
		for i := 0; i < len(layer.Strokes)-1; i++ {
			current := layer.Strokes[i]
			next := layer.Strokes[i+1]

			if current.Shape == Line && next.Shape == Line {
				// Draw line between consecutive line strokes
				rl.DrawLineEx(
					rl.Vector2{X: current.X, Y: current.Y},
					rl.Vector2{X: next.X, Y: next.Y},
					current.Size,
					current.Color,
				)

				// Draw endpoints
				rl.DrawCircle(int32(current.X), int32(current.Y), current.Size/2, current.Color)
				rl.DrawCircle(int32(next.X), int32(next.Y), current.Size/2, current.Color)

				// Skip the next stroke since we've processed it
				i++
			}
		}

		rl.EndTextureMode()
	}
}
