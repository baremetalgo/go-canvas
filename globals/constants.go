package globals

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

var CANVAS_WIDTH int32 = 800
var CANVAS_HEIGHT int32 = 600

type Stroke struct {
	Id         string
	X, Y       float32
	Size       float32
	Color      rl.Color
	Shape      int
	UsePattern bool
	Pattern    rl.Texture2D
}

func GenerateUniqueID() string {
	return uuid.New().String()
}
