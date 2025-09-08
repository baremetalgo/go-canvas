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
		texture.Format = rl.PixelFormat(8)
		BRUSH_PATTERNS = append(BRUSH_PATTERNS, texture)
		i += 1
		if i >= len(entries) {
			break
		}

	}
}

var AlphaDiscardShader rl.Shader

func Init() {
	AlphaDiscardShader = CreateAlphaDiscardShader()
}

// Alternatively, you can create the shader from code:
func CreateAlphaDiscardShader() rl.Shader {
	fsCode := `
    #version 330
    in vec2 fragTexCoord;
    in vec4 fragColor;
    out vec4 finalColor;
    uniform sampler2D texture0;
    void main() {
        vec4 texColor = texture(texture0, fragTexCoord);
        // Discard fragments with alpha below threshold
        if (texColor.a < 0.3) discard;
        finalColor = texColor * fragColor;
    }
    `

	shader := rl.LoadShaderFromMemory("", fsCode)
	return shader
}
