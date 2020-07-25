package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	//	"github.com/davecgh/go-spew/spew"
)

//var debug = true
var debug = false

// Process the input
type idx struct {
	X, Y,
	Wire int
}

// the value of the position indicates the number of steps until this point along the wire
type BoardB map[idx]int

var posX, posY, incrX, incrY = 0, 0, 0, 0
var maxX, minX, maxY, minY = 0, 0, 0, 0
var Wire, Pos = 0, 0

var board BoardB

func main() {
	board := make(BoardB)
	board[idx{0, 0, 0}] = 0

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		// /*
		if r := recover(); r != nil {
			fmt.Printf("w:%d p:%d \n", Wire, Pos)
			fmt.Printf("px:%d py:%d ix:%d iy:%d \n", posX, posY, incrX, incrY)
			fmt.Printf("X:%d x:%d Y:%d y:%d \n", maxX, minX, maxY, minY)
		}
		// */
	}()

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	// read the input from the specified file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}

	// convert the wire instructions into a local data structure
	for wire, line := range bytes.Split(data, []byte("\n")) {
		if wire > 1 {
			continue
		}
		posX, posY, incrX, incrY = 0, 0, 0, 0
		if debug {
			fmt.Printf("%d %+v \n", wire, string(line))
		}

		// we start at step one not zero since we've taken a step away from the origin
		step := 1
		for pos, elem := range bytes.Split(line, []byte(",")) {

			Wire = wire
			Pos = pos
			direction := string(elem[0])
			length, err := strconv.Atoi(string(elem[1:]))
			if debug {
				fmt.Printf("pos:%d elem:%v dir:%v l:%d err:%v\n", pos, string(elem), direction, length, err)
			}
			if err != nil {
				log.Panicf("Not a number:%+v", length, posX, posY, incrX, incrY)
			}
			switch direction {
			case "U":
				incrX, incrY = 0, 1
			case "D":
				incrX, incrY = 0, -1
			case "L":
				incrX, incrY = 1, 0
			case "R":
				incrX, incrY = -1, 0
			default:
				log.Panicf("Unexpected case: %s ", direction)
			}

			for n := 0; n < length; n++ {
				posX += incrX
				posY += incrY

				// keep track of our bounds
				if posX > maxX {
					maxX = posX
				}
				if posY > maxY {
					maxY = posY
				}
				if posX < minX {
					minX = posX
				}
				if posY < minY {
					minY = posY
				}
				board[idx{posX, posY, wire}] = step

				step++
			}
		}
	}

	if debug {
		fmt.Printf("w:%d p:%d \n", Wire, Pos)
		fmt.Printf("px:%d py:%d ix:%d iy:%d \n", posX, posY, incrX, incrY)
		fmt.Printf("X:%d x:%d Y:%d y:%d \n", maxX, minX, maxY, minY)
	}

	//	size := analyzeFullBoard(board, false)
	size := analyzeWires(board)

	fmt.Printf("Size: %d\n", size)
}

// output should be true if the board should be printed.
func analyzeFullBoard(board BoardB, output bool) int {
	type position struct{ X, Y, Distance int }
	var found []position

	for j := maxY; j >= minY; j-- {
		for i := maxX; i >= minX; i-- {
			//fmt.Printf(" %v|%v=%b=%b=",i,j,board[idx{i,j,0}],board[idx{i,j,1}])
			if i == j && j == 0 {
				// origin
				if output {
					fmt.Printf(" 000 ")
				}
				continue
			}
			if board[idx{i, j, 0}] != 0 && board[idx{i, j, 1}] != 0 {
				// crossover
				if output {
					fmt.Printf("-%3d-", board[idx{i, j, 0}]+board[idx{i, j, 1}])
				}
				// fmt.Sprintf("i:%d j:%d sum:%d\n", i, j, board[idx{i, j, 0}]+board[idx{i, j, 1}]
				found = append(found, position{X: i, Y: j, Distance: board[idx{i, j, 0}] + board[idx{i, j, 1}]})
				continue
			}
			if board[idx{i, j, 0}] != 0 {
				// wire 0 is present
				if output {
					fmt.Printf("*%3d*", board[idx{i, j, 0}])
				}
				continue
			}
			if board[idx{i, j, 1}] != 0 {
				// wire 1 is present
				if output {
					fmt.Printf("#%3d#", board[idx{i, j, 1}])
				}
				continue
			}
			if output {
				fmt.Printf("     ")
			}
		}
		if output {
			fmt.Print("\n")
		}
	}

	// now found is a set of intersection coordinates and the combined sum of the wires length
	// we need to range through the set and find the shortest length

	minDistance := 0
	for _, elem := range found {

		if elem.Distance < minDistance || minDistance == 0 {
			minDistance = elem.Distance
		}
	}

	fmt.Printf("found:\n%v\n", found)

	return minDistance
}

func analyzeWires(board BoardB) int {
	type position struct{ X, Y, Distance int }
	var found []position

	// this creates twice the work because we consider wire0 ?= wire1 and wire1 ?= wire0 but not twice the data
	// tw = this wire
	// ow = other wire (potentially)
	for tw, twDistance := range board {
		otherWire := 0
		if tw.Wire == 0 {
			otherWire = 1
		}

		if owDistance, ok := board[idx{X: tw.X, Y: tw.Y, Wire: otherWire}]; ok {
			// fmt.Println("-------------------------------------")
			// spew.Dump(i, distance, ow)
			found = append(found, position{X: tw.X, Y: tw.Y, Distance: twDistance + owDistance  } )
		}

	}

	// now found is a set of intersection coordinates and the combined sum of the wires length
	// we need to range through the set and find the shortest length

	minDistance := 0
	for _, elem := range found {

		if elem.Distance < minDistance || minDistance == 0 {
			minDistance = elem.Distance
		}
	}

	//	spew.Dump(found)
	
	return minDistance
}
