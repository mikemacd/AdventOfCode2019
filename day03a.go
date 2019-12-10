package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
)

//var debug = true
var debug = false
// Process the input
type idx struct{X,Y,Wire int}
type Board  map[idx]bool

var posX,posY,incrX,incrY = 0,0,0,0
var maxX,minX,maxY,minY = 0,0,0,0
var Wire, Pos = 0,0

var board Board
func main() {
	board := make(Board)
	board[idx{0,0,0}]=false

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		// /*
		if r := recover(); r != nil {
			fmt.Printf("w:%d p:%d \n", Wire, Pos)
			fmt.Printf("px:%d py:%d ix:%d iy:%d \n",posX,posY,incrX,incrY )
			fmt.Printf("X:%d x:%d Y:%d y:%d \n", maxX,minX,maxY,minY  )
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

	for wire, line := range bytes.Split(data, []byte("\n")) {
		if wire > 1 {
			continue
		}
		posX,posY,incrX,incrY = 0,0,0,0
		if debug {
			fmt.Printf("%d %+v \n", wire, string(line))
		}
		for pos, elem := range bytes.Split(line, []byte(",")) {
			Wire=wire
			Pos=pos
			direction := string(elem[0])
			length,err := strconv.Atoi(string(elem[1:]))
			if debug {
				fmt.Printf("pos:%d elem:%v dir:%v l:%d err:%v\n", pos, string(elem), direction, length, err)
			}
			if err!=nil {
				log.Panicf("Not a number:%+v",length, posX,posY,incrX,incrY)
			}
			switch direction {
			case "U":
				incrX,incrY =  0, 1
			case "D":
				incrX,incrY =  0, -1
			case "L":
				incrX,incrY =  1, 0
			case "R":
				incrX,incrY =  -1, 0
			default:
				log.Panicf("Unexpected case: %s ",direction)
			}

			for n:=0; n< length; n++ {
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
				board[idx{posX,posY,wire}]=true
			}
		}
	}

	if debug {
		fmt.Printf("w:%d p:%d \n", Wire, Pos)
		fmt.Printf("px:%d py:%d ix:%d iy:%d \n",posX,posY,incrX,incrY )
		fmt.Printf("X:%d x:%d Y:%d y:%d \n", maxX,minX,maxY,minY  )
	}

	// printBoard(board)

	size :=	analyzeBoard(board)

	fmt.Printf("Size: %d\n",size)
}

func printBoard(board Board){
	var found []string
	for j := maxY; j>=minY; j-- {
		for i:= maxX; i>=minX; i-- {
			//fmt.Printf(" %v|%v=%b=%b=",i,j,board[idx{i,j,0}],board[idx{i,j,1}])
			if i==j && j==0 {
				fmt.Printf("0")
				continue
			}
			if board[idx{i,j,0}]==true&&board[idx{i,j,1}]==true {
				fmt.Printf(":")
				found = append(found,fmt.Sprintf("i:%d j:%d\n",i,j))
				continue
			}
			if board[idx{i,j,0}]==true {
				fmt.Printf("^")
				continue
			}
			if board[idx{i,j,1}]==true {
				fmt.Printf(",")
				continue
			}
			fmt.Printf(" ")
		}
		fmt.Print("\n")
	}
	fmt.Printf("found:\n%v\n",found)
}

func analyzeBoard(board Board) (int) {
	var sizes []int
	for j := maxY; j>=minY; j-- {
		for i:= maxX; i>=minX; i-- {
			if board[idx{i, j, 0}] == true && board[idx{i, j, 1}] == true {
				currentSize := 0
				if debug {
					fmt.Printf("intersection found at i:%d j:%d \n",i,j)
				}
				if i<0 {
					currentSize += -i
				} else {
					currentSize += i
				}
				if j<0 {
					currentSize += -j
				} else {
					currentSize += j
				}

				sizes = append(sizes,currentSize)
			}
		}
	}
	minSize:=math.MaxInt64
	for _,v := range sizes {
		if  v<minSize {
			minSize=v
		}
	}
	if minSize < math.MaxInt64{
		return minSize
	}
	return 0
}