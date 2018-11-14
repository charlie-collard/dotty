package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/bspammer/dotty/braillify"
	"github.com/disintegration/imaging"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const defaultWidth int = 150
const defaultThreshold float64 = 0.5

var width = flag.Int("width", 0, "Width of the output text in characters")
var height = flag.Int("height", 0, "Height of the output text in characters")
var invert = flag.Bool("invert", false, "Invert the image")
var threshold = flag.Float64("threshold", defaultThreshold, "Level between 0 and 1 at which a pixel is considered bright enough to be a dot")
var choose = flag.Bool("choose", false, "Interactively choose a threshold to use")

func init() {
	flag.IntVar(width, "w", 0, "Width of the output text in characters")
	flag.IntVar(height, "h", 0, "Height of the output text in characters")
	flag.BoolVar(invert, "i", false, "Invert the image")
	flag.Float64Var(threshold, "t", defaultThreshold, "Level between 0 and 1 at which a pixel is considered bright enough to be a dot")
	flag.BoolVar(choose, "c", false, "Interactively choose a threshold to use")
}

func loadImage(filename string) (image.Image, error) {
	var imageStream io.ReadCloser
	// Fetch from the internet if the filename is a URL
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		response, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		imageStream = response.Body
	} else {
		var err error
		imageStream, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	defer imageStream.Close()

	img, _, err := image.Decode(imageStream)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	flag.Parse()
	if *width < 0 {
		log.Fatal(errors.New("dotty: width can't be negative"))
	}
	if *height < 0 {
		log.Fatal(errors.New("dotty: height can't be negative"))
	}
	if *threshold < 0 || *threshold > 1 {
		log.Fatal(errors.New("dotty: invalid threshold value"))
	}
	if flag.NArg() == 0 {
		log.Fatal(errors.New("dotty: missing filename for image"))
	}
	if *width == 0 && *height == 0 {
		*width = defaultWidth
	}
	// Braille characters are 2x4 grids of dots
	*width = *width * 2
	*height = *height * 4

	img, err := loadImage(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	resized := imaging.Resize(imaging.Grayscale(img), *width, *height, imaging.Lanczos)
	if *invert {
		resized = imaging.Invert(resized)
	}
	fmt.Println(braillify.ImgToBraille(resized, *threshold))
	if *choose {
		currentThreshold := *threshold
		reader := bufio.NewReader(os.Stdin)
	loop:
		for {
			currentThreshold = math.Max(currentThreshold, 0)
			currentThreshold = math.Min(currentThreshold, 1)
			fmt.Printf("[current threshold %.02f]\n", currentThreshold)
			fmt.Print("(u) up (U) 10x up (d) down (D) 10x down (s <val>) set (x) done: ")
			s, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			s = strings.Split(s, "\n")[0]
			switch s {
			case "u":
				currentThreshold += 0.01
			case "U":
				currentThreshold += 0.1
			case "d":
				currentThreshold -= 0.01
			case "D":
				currentThreshold -= 0.1
			case "x":
				break loop
			}
			if strings.HasPrefix(s, "s ") {
				s = strings.Split(s, " ")[1]
				currentThreshold, err = strconv.ParseFloat(s, 64)
				if err != nil {
					fmt.Println("not a valid threshold")
					continue
				}
			}
			fmt.Println(braillify.ImgToBraille(resized, currentThreshold))
		}
	}
}
