package widgets

import (
	"fmt"
	"go-canvas/globals"
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// FloodFill performs a flood fill operation on the specified layer
func FloodFill(layer *LayerSlot, startX, startY int, fillColor rl.Color, canvas *PaintCanvas) {
	// Convert RenderTexture2D to image for processing
	img := renderTextureToImage(layer.Texture)
	fmt.Println("sdsdsdsdsdsdsdsd")
	// Convert raylib color to Go image color
	targetColor := getColorAtPixel(img, startX, startY)
	fillImgColor := color.RGBA{
		R: fillColor.R,
		G: fillColor.G,
		B: fillColor.B,
		A: fillColor.A,
	}

	// Perform flood fill
	fill(img, startX, startY, targetColor, fillImgColor, 10000) // 10000 is max iterations

	// Convert image back to RenderTexture2D
	imageToRenderTexture(img, layer.Texture)
}

// Helper function to convert RenderTexture2D to image
func renderTextureToImage(tex rl.RenderTexture2D) *image.RGBA {
	// Read texture data
	rl.BeginTextureMode(tex)
	defer rl.EndTextureMode()

	// Create an image from the texture data
	// Note: This is a simplified approach - in a real implementation,
	// you'd need to properly read the texture data
	width := int(tex.Texture.Width)
	height := int(tex.Texture.Height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// For now, we'll return a blank image - you'd need to implement
	// proper texture data reading based on your specific setup
	return img
}

// Helper function to convert image back to RenderTexture2D
func imageToRenderTexture(img *image.RGBA, tex rl.RenderTexture2D) {
	// Implement texture updating based on your specific setup
	// This would involve updating the texture with the image data
}

// Get color at specific pixel
func getColorAtPixel(img *image.RGBA, x, y int) color.RGBA {
	if x < 0 || y < 0 || x >= img.Rect.Dx() || y >= img.Rect.Dy() {
		return color.RGBA{0, 0, 0, 0}
	}

	i := img.PixOffset(x, y)
	return color.RGBA{
		R: img.Pix[i],
		G: img.Pix[i+1],
		B: img.Pix[i+2],
		A: img.Pix[i+3],
	}
}

// Flood fill algorithm implementation
func fill(img *image.RGBA, x, y int, targetColor, fillColor color.RGBA, maxIterations int) {
	if targetColor.R == fillColor.R &&
		targetColor.G == fillColor.G &&
		targetColor.B == fillColor.B &&
		targetColor.A == fillColor.A {
		return // Target color same as fill color
	}

	queue := [][2]int{{x, y}}
	iterations := 0

	for len(queue) > 0 && iterations < maxIterations {
		iterations++
		point := queue[0]
		queue = queue[1:]

		x, y := point[0], point[1]

		// Check bounds
		if x < 0 || y < 0 || x >= img.Rect.Dx() || y >= img.Rect.Dy() {
			continue
		}

		// Get current pixel color
		currentColor := getColorAtPixel(img, x, y)

		// Check if we should fill this pixel
		if !colorsEqual(currentColor, targetColor) {
			continue
		}

		// Fill the pixel
		i := img.PixOffset(x, y)
		img.Pix[i] = fillColor.R
		img.Pix[i+1] = fillColor.G
		img.Pix[i+2] = fillColor.B
		img.Pix[i+3] = fillColor.A

		// Add neighbors to queue
		queue = append(queue,
			[2]int{x + 1, y}, // Right
			[2]int{x - 1, y}, // Left
			[2]int{x, y + 1}, // Down
			[2]int{x, y - 1}, // Up
		)
	}
}
func SimpleFloodFill(layer *LayerSlot, startX, startY int, fillColor rl.Color, canvas *PaintCanvas) {
	// Get the target color from the mouse position
	targetColor := getPixelColorFromTexture(layer.Texture, startX, startY)

	// If we're trying to fill with the same color, do nothing
	if colorsEqual(targetColor, fillColor) {
		return
	}

	// Create a temporary texture for the mask
	maskTex := rl.LoadRenderTexture(globals.CANVAS_WIDTH, globals.CANVAS_HEIGHT)
	defer rl.UnloadRenderTexture(maskTex)

	// Clear the mask
	rl.BeginTextureMode(maskTex)
	rl.ClearBackground(rl.Black)
	rl.EndTextureMode()

	// Create a mask of the area to fill using a recursive approach
	createFillMask(maskTex, layer.Texture, startX, startY, targetColor)

	// Apply the fill color to the masked area
	rl.BeginTextureMode(layer.Texture)
	rl.DrawTextureRec(
		maskTex.Texture,
		rl.NewRectangle(0, 0, float32(globals.CANVAS_WIDTH), -float32(globals.CANVAS_HEIGHT)),
		rl.NewVector2(0, 0),
		fillColor,
	)
	rl.EndTextureMode()
}

// createFillMask creates a mask of the area to fill
func createFillMask(maskTex, sourceTex rl.RenderTexture2D, x, y int, targetColor rl.Color) {
	// This is a simplified approach - in a real implementation,
	// you'd need to implement a proper flood fill algorithm

	// For now, we'll use a simple approach that fills a circular area
	// around the click point as a demonstration
	rl.BeginTextureMode(maskTex)
	rl.DrawCircle(int32(x), int32(y), 20, rl.White) // Fill a circular area
	rl.EndTextureMode()
}

// getPixelColorFromTexture gets the color of a pixel in a texture
func getPixelColorFromTexture(tex rl.RenderTexture2D, x, y int) rl.Color {
	// This is a placeholder - getting pixel colors from textures directly
	// is complex in raylib. In a real implementation, you'd need to:
	// 1. Read the texture data to CPU memory
	// 2. Extract the pixel color
	// 3. Process it

	// For demonstration, return a default color
	return rl.Black
}

// colorsEqual checks if two raylib colors are equal
func colorsEqual(c1, c2 rl.Color) bool {
	return c1.R == c2.R && c1.G == c2.G && c1.B == c2.B && c1.A == c2.A

}
