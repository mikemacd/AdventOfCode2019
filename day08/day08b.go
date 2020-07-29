package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type coord struct{ layer, x, y int }
type imageType struct {
	width  int
	height int
	layers int
	data   map[coord]int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}
	image := imageType{width: 25, height: 6, data: make(map[coord]int)}
	image.read(os.Args[1])

	_ = spew.Sdump(image.layers)
	//	fmt.Printf("Layer %d (which had %d zeros) has %d ones and %d twos with a product of %d\n", minZeroLayer, minZeroCount, ones, twos, result)
	image.render()
}

func (image *imageType) read(file string) {
	// raw reading of the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can't read file: %v\n", file)
		panic(err)
	}
	data = []byte(strings.TrimSpace(string(data)))

	layer := -1
	x := 0
	y := 0

	// take the read file and convert it from strings to ints
	for i, num := range bytes.Split(data, []byte{}) {
		switch {
		case i%(image.width*image.height) == 0:
			{
				layer++
				x = 0
				y = 0
			}
		case i%image.width == 0:
			{
				x = 0
				y++
			}
		default:
			x++
		}

		code, err := strconv.Atoi(string(num))
		if err != nil {
			log.Fatalf("Could not convert %v to integer. %v\n", num, err)
		}

		image.data[coord{layer, x, y}] = code

	}
	image.layers = layer
}

// return the number of zeros in a given layer
func (image *imageType) numDigits(digit, layer int) int {
	result := 0
	for j := 0; j < image.height; j++ {
		for i := 0; i < image.width; i++ {
			if image.data[coord{layer: layer, x: i, y: j}] == digit {
				result++
			}
		}
	}
	return result
}

func (image *imageType) print(layer int) {
	for j := 0; j < image.height; j++ {
		for i := 0; i < image.width; i++ {
			fmt.Printf("%v", image.data[coord{layer: layer, x: i, y: j}])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (image *imageType) render() {
	for j := 0; j < image.height; j++ {
		for i := 0; i < image.width; i++ {
			paint := " "

			for layer := image.layers; layer >= 0; layer-- {
				// for layer := 0 ; layer< image.layers; layer++ {
				pos := image.data[coord{layer: layer, x: i, y: j}]
				if pos != 2 {
					paint = fmt.Sprintf("%d", pos)

				}
			}
			if paint == "0" {
				fmt.Printf("  ")
			} else {
				fmt.Printf("\u2591\u2591")
			}

			//		fmt.Printf("%v", paint)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
