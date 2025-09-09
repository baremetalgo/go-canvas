package widgets

import (
	"fmt"
	"go-canvas/globals"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var BlendModes = map[string]rl.BlendMode{
	"BlendColors":     rl.BlendColor,
	"Alpha":           rl.BlendAlpha,
	"Additive":        rl.BlendAdditive,
	"AddColors":       rl.BlendAddColors,
	"BlendMultiply":   rl.BlendMultiplied,
	"SubtractColor":   rl.BlendSubtractColors,
	"PreMultAlpha":    rl.BlendAlphaPremultiply,
	"BlendDstAlpha":   rl.BlendDstAlpha,
	"BlendSrcAlpha":   rl.BlendSrcAlpha,
	"BlendSrcRGB":     rl.BlendSrcRgb,
	"BlenDstcRGB":     rl.BlendDstRgb,
	"BlendMultiplied": rl.BlendMultiplied,
}

type LayerSlot struct {
	Texture        rl.RenderTexture2D
	Strokes        []*globals.Stroke
	SlotID         float32
	SlotIndex      float32
	Editor         *LayerEditor
	Name           string
	Bounds         rl.Rectangle
	Visibility     bool
	DebugDraw      bool
	InvisibleColor rl.Color
	HighLightColor rl.Color

	BlendMode      string
	BlendToggle    bool
	Opacity        float32
	BlendTextColor rl.Color

	LayerNameBounds     rl.Vector2
	VisButton           rl.Rectangle
	BlendButton         rl.Rectangle
	deleteButton        rl.Rectangle
	delete_button_color rl.Color
	blendButtonColor    rl.Color
	BlendToggleButton   rl.Rectangle
	upButton            rl.Rectangle
	downButton          rl.Rectangle
	upcolor             rl.Color
	downcolor           rl.Color

	Active_blend_mode rl.Rectangle
	LeftButtonPressed bool

	delete_requested    bool
	move_up_requested   bool
	move_down_requested bool

	lastClickTime float64
	doubleClickMs float64
}

func NewLayerSlot(editor *LayerEditor, id float32) *LayerSlot {
	new_slot := LayerSlot{}
	new_slot.Texture = rl.LoadRenderTexture(globals.CANVAS_WIDTH, globals.CANVAS_HEIGHT)
	new_slot.SlotID = float32(id + 1)
	new_slot.SlotIndex = float32(id + 1)
	new_slot.Name = fmt.Sprintf("Layer_%v", id+1)
	new_slot.Visibility = true
	new_slot.Opacity = 1
	new_slot.BlendMode = "BlendColors"
	new_slot.BlendToggle = false
	new_slot.BlendTextColor = rl.DarkGreen
	new_slot.delete_button_color = rl.NewColor(255, 255, 255, 150)
	new_slot.blendButtonColor = rl.NewColor(10, 10, 10, 100)
	new_slot.upcolor = rl.NewColor(255, 255, 255, 150)
	new_slot.downcolor = rl.NewColor(255, 255, 255, 150)
	new_slot.Editor = editor
	new_slot.InvisibleColor = rl.NewColor(30, 30, 30, 75)
	new_slot.HighLightColor = rl.NewColor(255, 255, 255, 100)
	new_slot.Bounds = rl.NewRectangle(editor.Bounds.X, editor.Bounds.Y, editor.Bounds.Width, 50)
	new_slot.DebugDraw = false
	new_slot.Active_blend_mode = rl.NewRectangle(0, 0, 0, 0)
	new_slot.delete_requested = false
	new_slot.lastClickTime = 0
	new_slot.doubleClickMs = 300

	return &new_slot
}

func (s *LayerSlot) AddStroke(stroke *globals.Stroke) {
	s.Strokes = append(s.Strokes, stroke)
}

func (s *LayerSlot) RemoveStroke(stroke *globals.Stroke) {
	new_list := make([]*globals.Stroke, 0)
	for _, stroke_ := range s.Strokes {
		if stroke_.Id == stroke.Id {
			stroke_ = nil
			continue
		}
		new_list = append(new_list, stroke_)
	}
	s.Strokes = new_list
}

func (s *LayerSlot) GetStrokes() []*globals.Stroke {
	return s.Strokes
}

func (s *LayerSlot) GetId() float32 {
	return s.SlotID
}

func (s *LayerSlot) GetIndex() float32 {
	return s.SlotIndex
}

func (s *LayerSlot) GetBlendMode() string {
	return s.BlendMode
}

func (s *LayerSlot) GetOpacity() float32 {
	return s.Opacity
}

func (s *LayerSlot) GetTexture() rl.RenderTexture2D {
	return s.Texture
}
func (s *LayerSlot) GetVisibility() bool {
	return s.Visibility
}

func (s *LayerSlot) MakeActive() {
	s.Editor.ActiveLayer = s
}

func (s *LayerSlot) Clear() {
	// Clear strokes list
	s.Strokes = s.Strokes[:0]

	rl.BeginTextureMode(s.Texture)
	rl.ClearBackground(rl.Blank)
	rl.EndTextureMode()
}

func (s *LayerSlot) Update() {
	s.move_up_requested = false
	s.move_down_requested = false
	s.Bounds = rl.NewRectangle(s.Editor.Bounds.X, s.Editor.Bounds.Y, s.Editor.Bounds.Width, 50)

	if s.SlotIndex == 1 {
		s.Bounds.Y = s.Bounds.Y + float32(s.Editor.TitleBarHeight) + 15
	} else {
		s.Bounds.Y = s.Bounds.Y + (s.SlotIndex * s.Bounds.Height) + (15 * s.SlotIndex) - float32(s.Editor.TitleBarHeight)
	}
	s.VisButton = rl.NewRectangle(s.Bounds.X+5, s.Bounds.Y+(s.Bounds.Height)/2-10, 20, 20)
	s.LayerNameBounds = rl.NewVector2(s.VisButton.X+25, s.VisButton.Y+5)
	layer_name_length := float32(rl.MeasureText(s.Name, 14))
	s.BlendButton = rl.NewRectangle(s.LayerNameBounds.X+layer_name_length+25, s.LayerNameBounds.Y, 50, 20)
	s.BlendToggleButton = rl.NewRectangle(s.LayerNameBounds.X+layer_name_length, s.LayerNameBounds.Y-5, 20, 20)
	s.deleteButton = rl.NewRectangle(s.Bounds.X+s.Bounds.Width-25, s.Bounds.Y+s.Bounds.Height/2-10, 20, 20)
	s.upButton = rl.NewRectangle(s.Bounds.X+s.Bounds.Width, s.Bounds.Y+10, 15, 15)
	s.downButton = rl.NewRectangle(s.Bounds.X+s.Bounds.Width, s.Bounds.Y+s.Bounds.Height-20, 15, 15)

	mouse_pos := rl.GetMousePosition()

	// Middle mouse drag for opacity control
	if rl.CheckCollisionPointRec(mouse_pos, s.Bounds) {
		rl.DrawRectangleLinesEx(s.Bounds, 2, rl.White)

		// Middle mouse button drag for opacity
		if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
			// Get mouse movement delta
			mouse_delta := rl.GetMouseDelta()

			// Calculate opacity change based on horizontal movement
			opacity_change := mouse_delta.X * 0.01 // Adjust sensitivity as needed

			// Update opacity, clamping between 0 and 1
			s.Opacity += opacity_change
			if s.Opacity < 0 {
				s.Opacity = 0
			} else if s.Opacity > 1 {
				s.Opacity = 1
			}
		}

		// Double click to make active
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			now := rl.GetTime() * 1000 // ms
			if now-s.lastClickTime <= s.doubleClickMs {
				// Double click detected
				s.MakeActive()
				fmt.Printf("Layer %s (ID %.0f) set active\n", s.Name, s.SlotID)
			}
			s.lastClickTime = now
		}
	}

	if rl.CheckCollisionPointRec(mouse_pos, s.VisButton) {
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			s.Visibility = !s.Visibility
		}
	}

	if rl.CheckCollisionPointRec(mouse_pos, s.BlendButton) {
		s.BlendTextColor = rl.Red
	} else {
		s.BlendTextColor = rl.DarkGreen
	}
	if rl.CheckCollisionPointRec(mouse_pos, s.BlendButton) && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		blendModeKeys := []string{"BlendColors", "Alpha", "Additive", "AddColors", "BlendMultiply", "SubtractColor",
			"PreMultAlpha", "BlendDstAlpha", "BlendSrcAlpha", "BlendSrcRGB", "BlenDstcRGB", "BlendMultiplied"}
		next_blend_mode := nextElement(blendModeKeys, s.BlendMode)
		s.BlendMode = next_blend_mode
	}

	if rl.CheckCollisionPointRec(mouse_pos, s.deleteButton) {
		s.delete_button_color = rl.NewColor(255, 255, 255, 255)
	} else {
		s.delete_button_color = rl.NewColor(255, 255, 255, 150)
	}

	if rl.CheckCollisionPointRec(mouse_pos, s.deleteButton) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		s.delete_requested = true

	}

	if rl.CheckCollisionPointRec(mouse_pos, s.upButton) {
		s.upcolor = rl.NewColor(255, 255, 255, 255)
	} else {
		s.upcolor = rl.NewColor(255, 255, 255, 100)
	}
	if rl.CheckCollisionPointRec(mouse_pos, s.upButton) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		s.move_up_requested = true
	}

	if rl.CheckCollisionPointRec(mouse_pos, s.downButton) {
		s.downcolor = rl.NewColor(255, 255, 255, 255)
	} else {
		s.downcolor = rl.NewColor(255, 255, 255, 100)
	}
	if rl.CheckCollisionPointRec(mouse_pos, s.downButton) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		s.move_down_requested = true
	}
	if rl.CheckCollisionPointRec(mouse_pos, s.BlendToggleButton) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		s.BlendToggle = !s.BlendToggle
		if s.BlendToggle {
			s.blendButtonColor = rl.White
		} else {
			s.blendButtonColor = rl.NewColor(10, 10, 10, 100)
		}
	}
}

func (s *LayerSlot) Draw() {
	s.Update()

	// draw body
	if s.Visibility {
		rl.DrawRectangle(
			s.Bounds.ToInt32().X,
			s.Bounds.ToInt32().Y,
			s.Bounds.ToInt32().Width,
			s.Bounds.ToInt32().Height,
			rl.NewColor(50, 50, 50, 30))
	} else {
		rl.DrawRectangle(
			s.Bounds.ToInt32().X,
			s.Bounds.ToInt32().Y,
			s.Bounds.ToInt32().Width,
			s.Bounds.ToInt32().Height,
			s.InvisibleColor)
	}
	// draw visibility round

	vis_texture := SlotIcons["eye_open"]
	if !s.Visibility {
		vis_texture = SlotIcons["eye_close"]
	}

	src := rl.NewRectangle(0, 0, float32(vis_texture.Width), float32(vis_texture.Height))
	rl.DrawTexturePro(vis_texture, src, s.VisButton, rl.NewVector2(0, 0), 0, rl.White)

	// Draw layer name
	new_name := s.StringPadRightBuilder(s.Name, 10)
	s.Name = new_name
	rl.DrawTextPro(Default_Widget_Header_Font, new_name, s.LayerNameBounds, rl.NewVector2(0, 0), 0, 14, 0, rl.Black)

	// Draw Blend Icon
	icon_texture := SlotIcons["blend_icon"]
	src = rl.NewRectangle(0, 0, float32(icon_texture.Width), float32(icon_texture.Height))
	rl.DrawTexturePro(icon_texture, src, s.BlendToggleButton, rl.NewVector2(0, 0), 0, s.blendButtonColor)

	// Draw Delete button
	delete_texture := SlotIcons["delete_icon"]
	src = rl.NewRectangle(0, 0, float32(delete_texture.Width), float32(delete_texture.Height))
	// dst = rl.NewRectangle(s.deleteButton.X, s.deleteButton.Y+s.Bounds.Height/2-10, 20, 20)
	rl.DrawTexturePro(delete_texture, src, s.deleteButton, rl.NewVector2(0, 0), 0, s.delete_button_color)

	// Draw Blend mode text
	rl.DrawTextEx(Default_Widget_Header_Font, s.BlendMode, rl.NewVector2(s.BlendButton.X, s.BlendButton.Y), 14, 2, s.BlendTextColor)

	// Drawing opacity overlay
	opacity_rect := rl.NewRectangle(s.Bounds.X, s.Bounds.Y, s.Bounds.Width*s.Opacity, s.Bounds.Height)
	rl.DrawRectanglePro(opacity_rect, rl.NewVector2(0, 0), 0, rl.NewColor(0, 0, 0, 20))

	// draw up buttons
	up_text := SlotIcons["up_arrow"]
	src = rl.NewRectangle(0, 0, float32(up_text.Width), float32(up_text.Height))
	rl.DrawTexturePro(up_text, src, s.upButton, rl.NewVector2(0, 0), 0, s.upcolor)
	// draw down buttons
	down_text := SlotIcons["down_arrow"]
	src = rl.NewRectangle(0, 0, float32(down_text.Width), float32(down_text.Height))
	rl.DrawTexturePro(down_text, src, s.downButton, rl.NewVector2(0, 0), 0, s.downcolor)

	// Draw debug
	if s.DebugDraw {
		rl.DrawRectangleLinesEx(s.VisButton, 2, rl.Red)
		rl.DrawRectangleLinesEx(s.Bounds, 2, rl.Red)
	}
}

func (slot *LayerSlot) StringPadRightBuilder(s string, length int) string {
	if len(s) >= length {
		return s
	}

	var sb strings.Builder
	sb.WriteString(s)
	for i := 0; i < length-len(s); i++ {
		sb.WriteRune(' ')
	}
	return sb.String()
}

func nextElement(slice []string, current string) string {
	for i, v := range slice {
		if v == current {
			// if last element, wrap to first
			if i == len(slice)-1 {
				return slice[0]
			}
			// else return next element
			return slice[i+1]
		}
	}
	return "" // not found
}
