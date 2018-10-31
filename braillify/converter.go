// Package braillify implements conversion from images to braille strings
package braillify

import (
	"image"
	"strings"
)

// ImgToBraille produces a string of braille unicode characters representing an image.
// Each pixel in the input image is represented as a single dot in a braille character in the output.
// Braille characters are 2x4, so a 200x100 image would be converted into 100x25 braille characters.
// threshold is a value between 0 and 1 representing the level of brightness needed to display a dot
func ImgToBraille(img image.Image, threshold float64) string {
	if threshold > 1 {
		threshold = 1
	}
	if threshold < 0 {
		threshold = 0
	}
	var extra int
	if img.Bounds().Dy()%4 != 0 {
		extra = 1
	}
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()
	brailleWidth := int(imageWidth / 2)
	brailleHeight := int(extra + imageHeight/4)
	brailleArray := make([]uint8, brailleWidth*brailleHeight)

	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			grayscale := (r + g + b) / 3
			if grayscale >= uint32(threshold*0xffff) {
				brailleArray[int(i/2)+int(j/4)*brailleWidth] |= 1 << uint8((i&1)*4+(j&3))
			}
		}
	}
	var b strings.Builder
	for i, c := range brailleArray {
		b.WriteRune(brailleMap[c])
		if (i%brailleWidth) == 0 && i != 0 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}
