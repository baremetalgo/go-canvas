package globals

import (
	"fmt"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

var CANVAS_WIDTH int32 = 800
var CANVAS_HEIGHT int32 = 600
var BRUSH_PATTERNS []rl.Texture2D

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

func LoadToolBoxIcons() {
	entries, err := os.ReadDir("./patterns")
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for i < 12 {
		texture := rl.LoadTexture(fmt.Sprintf("patterns/%v", entries[i].Name()))
		BRUSH_PATTERNS = append(BRUSH_PATTERNS, texture)
		i += 1
		if i >= len(entries) {
			break
		}

	}
}
