package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Layer struct {
	Visibility bool
	Texture    rl.RenderTexture2D
	Id         uint32
	Name       string
	BlendMode  rl.BlendMode
	Opacity    int32
	Strokes    []Stroke
}

func NewLayer() *Layer {
	canvas := rl.LoadRenderTexture(800, 600)
	return &Layer{
		Visibility: true,
		Texture:    canvas,
		Id:         canvas.ID,
		Name:       fmt.Sprintf("Layer_%v", canvas.ID),
		BlendMode:  rl.BlendAddColors,
		Opacity:    255,
	}
}
