// Package braillify implements conversion from images to braille strings
package braillify

import (
	"image"
	"math"
	"strings"
)

// ImgToBraille produces a string of braille unicode characters representing an image.
// Each pixel in the input image is represented as a single dot in a braille character in the output.
// Braille characters are 2x4, so a 200x100 image would be converted into 100x25 braille characters.
// threshold is a value between 0 and 1 representing the level of brightness needed to display a dot
func ImgToBraille(img image.Image, threshold float64) string {
	// Clamp the possible threshold values
	threshold = math.Max(threshold, 0)
	threshold = math.Min(threshold, 1)
	var extra int
	// Add an extra row if the height of the image is not exactly divisible by 4
	if img.Bounds().Dy()%4 != 0 {
		extra = 1
	}
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()
	brailleWidth := int(imageWidth/2) + imageWidth%2
	brailleHeight := int(imageHeight/4) + extra
	brailleArray := make([]uint8, brailleWidth*brailleHeight)

	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
			x, y := i+img.Bounds().Min.X, j+img.Bounds().Min.Y
			r, g, b, _ := img.At(x, y).RGBA()
			grayscale := (r + g + b) / 3
			if grayscale > uint32(threshold*0xffff) {
				brailleArray[int(i/2)+int(j/4)*brailleWidth] |= 1 << uint8((i&1)*4+(j&3))
			}
		}
	}
	var b strings.Builder
	for i, c := range brailleArray {
		if (i%brailleWidth) == 0 && i != 0 {
			b.WriteRune('\n')
		}
		b.WriteRune(brailleMap[c])
	}
	return b.String()
}
