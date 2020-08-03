package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type pos struct {
	X int
	Y int
}
type chartType struct {
	width     int
	height    int
	positions map[pos]string
	asteroids []pos
}

var debug = false

func main() {
	_ = spew.Sdump("")

	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}

	chart := readChart(os.Args[1])

	chart.print()

	num, pos := chart.bestLocation()

	fmt.Printf("\n%d,%d can see %d other asteroids\n\n", pos.X, pos.Y, num)

	chart.printHighlight(pos)

	fmt.Println("")

	chart.printWithVisible()
}

func readChart(file string) chartType {
	chart := chartType{positions: map[pos]string{}, asteroids: []pos{}}

	// raw reading of the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can't read file: %v\n", file)
		panic(err)
	}

	for j, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if len(line) > chart.width {
			chart.width = len(line)
		}
		chart.height++
		for i := 0; i < len(line); i++ {
			chart.positions[pos{X: i, Y: j}] = string(line[i])
			if line[i] == '#' {
				chart.asteroids = append(chart.asteroids, pos{X: i, Y: j})
			}
		}
	}

	return chart
}

func (c *chartType) print() {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			fmt.Printf("%s", c.positions[pos{X: i, Y: j}])
		}
		fmt.Print("\n")
	}
}

func (c *chartType) printWithVisible() {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			if c.positions[pos{X: i, Y: j}] != "#" {
				fmt.Print("... ")
				continue
			}
			vAst := c.numberOfAsteroidsVisible(pos{X: i, Y: j})
			fmt.Printf("%3d ", len(vAst))
		}
		fmt.Print("\n")
	}
}

func (c *chartType) printHighlight(p pos) {
	for j := 0; j < c.height; j++ {
		for i := 0; i < c.width; i++ {
			char := c.positions[pos{X: i, Y: j}]
			if p.X == i && p.Y == j {
				char = "\u2588"
			}
			fmt.Printf("%s", char)
		}
		fmt.Print("\n")
	}
}

func (c *chartType) bestLocation() (int, pos) {
	maxV := 0
	maxP := pos{}
	allVis := ""
	for _, p := range c.asteroids {
		vAst := c.numberOfAsteroidsVisible(p)
		v := len(vAst)
		allVis += fmt.Sprintf("Visible asteroids to %v :: %v\n", p, vAst)
		if v > maxV {
			maxV = v
			maxP = p
		}
	}
	if debug{
		fmt.Print(allVis)	
	}

	return maxV, maxP
}

var reasons = ""

func (c *chartType) numberOfAsteroidsVisible(astA pos) []pos {
	// asteroidA can see asteroidC (an asteroid from set c.asteroid that isnt asteroidA)
	// if there is no asteroidB (an asterpoid from set c.asteroid that isn't asteroidC or asteroidA)
	// such that slope AB = slope BC
	// if asteroidA can see asteroidC add it to the collection (array) of asteroids that can be seen by asteroidC
	// what is the size of the collection?

	// visible== asteroidCs
	// for all asteroidCs , and for all other asteroidBs if |AB| == |AC| && A->B->C pop asteroid from visible 

 	// all asteroids could be able to be seen by astA except:
	visibleAsts := make([]pos, len(c.asteroids))
	copy(visibleAsts, c.asteroids)
	// if they _are_ asteroidA
	visibleAsts = pop(visibleAsts, astA)

	for _, astC := range visibleAsts {
		// an other asteroid is an asteroid that can be seen by asteroid asteroid
		otherAsts := make([]pos, len(visibleAsts))
		copy(otherAsts, visibleAsts)
		// that isn't asteroidC
		otherAsts = pop(otherAsts, astC)

		// for all other
		for _, astB := range otherAsts {
 
			if collinear(astA,astB,astC) && ordered(astA,astB,astC) {	
				visibleAsts = pop(visibleAsts, astC)
			}
		}
	}
	return visibleAsts
}

func pop(set []pos, elem pos) []pos {
	rv := []pos{}
	for _, v := range set {
		if v != elem {
			rv = append(rv, v)
		}
	}
	return rv
}

// are the points collinear A->C->B or B->C->A or A->B->C etc ?
func collinear(astA,astB,astC pos) (bool) {
	return 0 == (astB.Y-astA.Y)*(astC.X-astB.X)-(astB.X-astA.X)*(astC.Y-astB.Y)
}

// are they orderd A->B->C ?
func ordered(astA,astB,astC pos) (bool) {
	return ((astA.X <= astB.X && astB.X <= astC.X) || (astA.X >= astB.X && astB.X >= astC.X)) && ((astA.Y <= astB.Y && astB.Y <= astC.Y) || (astA.Y >= astB.Y && astB.Y >= astC.Y))
}
