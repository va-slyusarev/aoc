package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"unicode"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var verbose = flag.Bool("v", false, "Verbose output")

var intersectionColor color.Color = color.RGBA{R: 0xff, A: 0xff}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines [][]string

	for scanner.Scan() {
		lines = append(lines, strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			return !unicode.IsDigit(r) && !unicode.IsLetter(r)
		}))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	img := paveWires(rect(lines), lines)

	if *verbose {
		imgFile, err := os.Create("image.png")
		defer func() { _ = imgFile.Close() }()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		if err := png.Encode(imgFile, img); err != nil {
			log.Fatalf("%v\n", err)
		}
	}

	fmt.Printf("Success! Target number is: %d\n", manhattanDist(img))
}

func rect(lines [][]string) image.Rectangle {
	min := image.Pt(0, 0)
	max := image.Pt(0, 0)

	for _, line := range lines {

		current := image.Pt(0, 0)

		for _, direction := range line {
			var dir string
			var value int
			if _, err := fmt.Sscanf(direction, "%1s%d", &dir, &value); err != nil {
				log.Fatalf("Error parsing direction: %v\n", err)
			}
			switch dir {
			case "U":
				current = current.Add(image.Pt(0, value))
				if max.Y < current.Y {
					max.Y = current.Y
				}
			case "R":
				current = current.Add(image.Pt(value, 0))
				if max.X < current.X {
					max.X = current.X
				}
			case "D":
				current = current.Add(image.Pt(0, -value))
				if min.Y > current.Y {
					min.Y = current.Y
				}
			case "L":
				current = current.Add(image.Pt(-value, 0))
				if min.X > current.X {
					min.X = current.X
				}
			default:
				log.Fatalf("Unknown direction: %v\n", dir)
			}
		}
	}

	// extended area
	min = min.Sub(image.Pt(1, 1))
	max = max.Add(image.Pt(2, 2))

	return image.Rectangle{Min: min, Max: max}
}

func paveWires(rect image.Rectangle, lines [][]string) *image.RGBA {

	img := image.NewRGBA(rect)

	for idx, line := range lines {

		current := image.Pt(0, 0)
		col := color.RGBA{R: uint8(idx + 1), G: uint8(200 * idx), B: uint8(100), A: 0xff}

		for _, direction := range line {
			var dir string
			var value int
			if _, err := fmt.Sscanf(direction, "%1s%d", &dir, &value); err != nil {
				log.Fatalf("Error parsing direction: %v\n", err)
			}
			switch dir {
			case "U":
				for y := 1; y <= value; y++ {
					dx, dy := current.X, current.Y+y
					if img.RGBAAt(dx, dy).R != 0 && img.RGBAAt(dx, dy).R != uint8(idx+1) {
						img.Set(dx, dy, intersectionColor)
					} else {
						img.Set(dx, dy, col)
					}
				}
				current = current.Add(image.Pt(0, value))

			case "R":
				for x := 1; x <= value; x++ {
					dx, dy := current.X+x, current.Y
					if img.RGBAAt(dx, dy).R != 0 && img.RGBAAt(dx, dy).R != uint8(idx+1) {
						img.Set(dx, dy, intersectionColor)
					} else {
						img.Set(dx, dy, col)
					}
				}
				current = current.Add(image.Pt(value, 0))
			case "D":
				for y := 1; y <= value; y++ {
					dx, dy := current.X, current.Y-y
					if img.RGBAAt(dx, dy).R != 0 && img.RGBAAt(dx, dy).R != uint8(idx+1) {
						img.Set(dx, dy, intersectionColor)
					} else {
						img.Set(dx, dy, col)
					}
				}
				current = current.Add(image.Pt(0, -value))
			case "L":
				for x := 1; x <= value; x++ {
					dx, dy := current.X-x, current.Y
					if img.RGBAAt(dx, dy).R != 0 && img.RGBAAt(dx, dy).R != uint8(idx+1) {
						img.Set(dx, dy, intersectionColor)
					} else {
						img.Set(dx, dy, col)
					}
				}
				current = current.Add(image.Pt(-value, 0))
			default:
				log.Fatalf("Unknown direction: %v\n", dir)
			}
		}
	}
	img.Set(0, 0, color.White)
	return img
}

func manhattanDist(img *image.RGBA) int {
	result := math.MaxInt32
	for x := img.Rect.Min.X; x <= img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y <= img.Rect.Max.Y; y++ {

			// intersection wire's
			if img.At(x, y) != intersectionColor {
				continue
			}

			dist := int(math.Abs(float64(x)) + math.Abs(float64(y)))
			if dist < result {
				result = dist
			}
		}
	}
	return result
}
