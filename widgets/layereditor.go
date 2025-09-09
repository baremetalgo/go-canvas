package widgets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SlotIcons map[string]rl.Texture2D

type LayerEditor struct {
	BaseWidget
	Slots           []*LayerSlot
	ActiveLayer     *LayerSlot
	BrushColor      rl.Color
	HightlightColor rl.Color
	AddLayerButton  rl.Rectangle
	DebugDraw       bool
}

func NewLayerEditor(name string) *LayerEditor {
	editor := LayerEditor{}
	editor.LoadSlotIcons()
	// Setup BaseWidget first
	editor.Slots = make([]*LayerSlot, 0)
	editor.BrushColor = rl.White
	editor.Visible = true
	editor.Name = name
	editor.DrawTitleBar = true
	editor.TitleBarHeight = 20
	editor.TitleBarColor = rl.DarkGray
	editor.TextColor = rl.White
	editor.BodyColor = rl.NewColor(50, 50, 50, 0)
	editor.BorderColor = rl.NewColor(50, 50, 50, 20)
	editor.BaseWidget.BodyColor = editor.BodyColor
	editor.DrawBorder = false
	editor.Width = 250
	editor.Height = 400
	editor.Bounds = rl.NewRectangle(float32(rl.GetScreenWidth()-int(editor.Width)-15), 50, float32(editor.Width), float32(editor.Height))
	editor.TitleBarBounds = rl.NewRectangle(editor.Bounds.X, editor.Bounds.Y, float32(editor.Width), float32(editor.TitleBarHeight))
	editor.AddLayerButton = rl.NewRectangle(editor.Bounds.X, editor.Bounds.Y, 20, 20)
	editor.HightlightColor = rl.NewColor(255, 255, 255, 50)
	editor.DebugDraw = false
	editor.AddSlot()
	editor.ActiveLayer = editor.Slots[0]
	return &editor
}

func (e *LayerEditor) Draw() {
	e.BaseWidget.Draw()
	e.Update()
	for _, slot := range e.Slots {
		slot.Draw()
	}

	// Draw Add Layer button
	tex := SlotIcons["add"]
	src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
	rl.DrawTexturePro(tex, src, e.AddLayerButton, rl.NewVector2(0, 0), 0, rl.White)

	// Draw active Layer highlight
	if e.ActiveLayer != nil {
		rl.DrawRectangleLinesEx(e.ActiveLayer.Bounds, 3, rl.White)
	}
	e.DeleteSlots()
	e.MoveSlots()
}

func (e *LayerEditor) Update() {
	total_slot_height := len(e.Slots)*70 + int(e.TitleBarHeight)
	e.Height = int32(total_slot_height)
	e.BaseWidget.Height = e.Height

	if len(e.Slots) > 0 {
		last_slot := e.Slots[len(e.Slots)-1]
		last_slot_y_pos := last_slot.Bounds.Y + last_slot.Bounds.Height
		e.AddLayerButton.X = e.Bounds.X
		e.AddLayerButton.Y = last_slot_y_pos
	} else {
		e.AddLayerButton.X = e.Bounds.X
		e.AddLayerButton.Y = float32(e.Height) + 30 + float32(e.TitleBarHeight)
	}

	mouse_pos := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mouse_pos, e.AddLayerButton) {
		rl.DrawCircle(e.AddLayerButton.ToInt32().X+10, e.AddLayerButton.ToInt32().Y+10, 15, rl.White)
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			e.AddSlot()
		}
	}
}

func (e *LayerEditor) AddSlot() {
	new_slot := NewLayerSlot(e, float32(len(e.Slots)))
	e.Slots = append(e.Slots, new_slot)
	e.ActiveLayer = new_slot

}

func (e *LayerEditor) RemoveSlot(SlotID float32) {
	new_slots := make([]*LayerSlot, 0)
	for _, slot := range e.Slots {
		if slot.SlotID != SlotID {
			new_slots = append(new_slots, slot)
			slot.SlotIndex = float32(len(new_slots))
		} else {
			slot = nil
		}

	}
	e.Slots = new_slots

}

func (e *LayerEditor) MoveSlots() {
	for i := 0; i < len(e.Slots); i++ {
		slot := e.Slots[i]

		// Move Up
		if slot.move_up_requested && i > 0 {
			// swap with previous
			e.Slots[i], e.Slots[i-1] = e.Slots[i-1], e.Slots[i]

			// reset flag
			slot.move_up_requested = false
		}

		// Move Down
		if slot.move_down_requested && i < len(e.Slots)-1 {
			// swap with next
			e.Slots[i], e.Slots[i+1] = e.Slots[i+1], e.Slots[i]

			// reset flag
			slot.move_down_requested = false
			// since we swapped with next, skip the next index to avoid double-processing
			i++
		}
	}

	// Re-index all slots
	for idx, slot := range e.Slots {
		slot.SlotIndex = float32(idx + 1)
	}
}

func (e *LayerEditor) DeleteSlots() {
	for _, slot := range e.Slots {
		if slot.delete_requested {
			e.RemoveSlot(slot.SlotID)
		}
	}
}

func (e *LayerEditor) LoadSlotIcons() {
	SlotIcons = map[string]rl.Texture2D{
		"eye_open":    rl.LoadTexture("icons/visible_icon.png"),
		"eye_close":   rl.LoadTexture("icons/invisible_icon.png"),
		"blend_icon":  rl.LoadTexture("icons/layers_icon.png"),
		"delete_icon": rl.LoadTexture("icons/bin_icon.png"),
		"up_arrow":    rl.LoadTexture("icons/up_arrow.png"),
		"down_arrow":  rl.LoadTexture("icons/down_arrow.png"),
		"add":         rl.LoadTexture("icons/add_icon.png"),
	}
}
