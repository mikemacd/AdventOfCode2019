package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type pos struct {
	X int
	Y int
}

type posSet []pos

type orbitSet struct {
	origin pos
	set    posSet
}

type chartType struct {
	width     int
	height    int
	positions map[pos]string
	asteroids posSet
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

	num, p := chart.bestLocation()

	fmt.Printf("\n%d,%d can see %d other asteroids\n\n", p.X, p.Y, num)

	chart.printHighlight(p)

	fmt.Println("")

	vAstSet := chart.numberOfAsteroidsVisible(p)
	vAstPositions := map[pos]string{}

	for j := 0; j < chart.height; j++ {
		for i := 0; i < chart.width; i++ {
			vAstPositions[pos{X: i, Y: j}] = "."
		}
	}

	for _, v := range chart.asteroids {
		for _, q := range vAstSet {
			if q == v {
				vAstPositions[q] = "@"
			}
		}
	}
	fmt.Printf("AAA %d %d\n", len(vAstSet), len(chart.asteroids))
	// chart.printWithVisible()
	vAstChart := chartType{
		width:     chart.width,
		height:    chart.height,
		positions: vAstPositions,
		asteroids: vAstSet,
	}
	vAstChart.printHighlight(p)
	vAst := orbitSet{
		set:    vAstSet,
		origin: p,
	}
	fmt.Println()
	fmt.Printf("set before sorting\n")

	for i, v := range vAst.set {
		fmt.Printf("\t%d:%d", i, v)
	}

	fmt.Println()
	fmt.Println()

	sort.Sort(vAst)
	fmt.Printf("#############\n")
	fmt.Println()
	fmt.Printf("set after sorting\n")

	for i, v := range vAst.set {
		fmt.Printf("\t%d:%d", i, v)
		if debug {
			vecA := pos{
				X: (v.X - vAst.origin.X),
				Y: (v.Y - vAst.origin.Y),
			}

			angleA := angleWithTwelveOclock(vecA)
			fmt.Printf("%v:%v\tO:%v a:%d  vecA:%v   angleA:%2.4f \t%v \n", i, v, vAst.origin, v, vecA, angleA, angleA*180/math.Pi)
		}
	}
	// get the 200th asteroid from the 199th position
	rv := vAst.set[199].X*100 + vAst.set[199].Y

	fmt.Printf("\n\nThe encoded 200th position is: %d\n", rv)
}

func readChart(file string) chartType {
	chart := chartType{positions: map[pos]string{}, asteroids: posSet{}}

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
	fmt.Print("    ")
	k := 0
	for i := 0; i < c.width; i++ {
		k++
		if (k-1)%10 == 0 {
			fmt.Printf("%d", (i / 10))
			continue
		}
		fmt.Print(" ")
	}
	fmt.Println()
	fmt.Println()
	fmt.Print("    ")
	for i := 0; i < c.width; i++ {
		fmt.Printf("%d", i%10)
	}
	fmt.Println()
	for j := 0; j < c.height; j++ {
		fmt.Printf("%2d  ", j)
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
	if debug {
		fmt.Print(allVis)
	}

	return maxV, maxP
}

var reasons = ""

func (c *chartType) numberOfAsteroidsVisible(astA pos) posSet {
	// asteroidA can see asteroidC (an asteroid from set c.asteroid that isnt asteroidA)
	// if there is no asteroidB (an asterpoid from set c.asteroid that isn't asteroidC or asteroidA)
	// such that slope AB = slope BC
	// if asteroidA can see asteroidC add it to the collection (array) of asteroids that can be seen by asteroidC
	// what is the size of the collection?

	// visible== asteroidCs
	// for all asteroidCs , and for all other asteroidBs if |AB| == |AC| && A->B->C pop asteroid from visible

	// all asteroids could be able to be seen by astA except:
	visibleAsts := make(posSet, len(c.asteroids))
	copy(visibleAsts, c.asteroids)
	// if they _are_ asteroidA
	visibleAsts = pop(visibleAsts, astA)

	for _, astC := range visibleAsts {
		// an other asteroid is an asteroid that can be seen by asteroid asteroid
		otherAsts := make(posSet, len(visibleAsts))
		copy(otherAsts, visibleAsts)
		// that isn't asteroidC
		otherAsts = pop(otherAsts, astC)

		// for all other
		for _, astB := range otherAsts {

			if collinear(astA, astB, astC) && ordered(astA, astB, astC) {
				visibleAsts = pop(visibleAsts, astC)
			}
		}
	}
	return visibleAsts
}

func pop(set posSet, elem pos) posSet {
	rv := posSet{}
	for _, v := range set {
		if v != elem {
			rv = append(rv, v)
		}
	}
	return rv
}

// are the points collinear A->C->B or B->C->A or A->B->C etc ?
func collinear(astA, astB, astC pos) bool {
	return 0 == (astB.Y-astA.Y)*(astC.X-astB.X)-(astB.X-astA.X)*(astC.Y-astB.Y)
}

// are they orderd A->B->C ?
func ordered(astA, astB, astC pos) bool {
	return ((astA.X <= astB.X && astB.X <= astC.X) || (astA.X >= astB.X && astB.X >= astC.X)) && ((astA.Y <= astB.Y && astB.Y <= astC.Y) || (astA.Y >= astB.Y && astB.Y >= astC.Y))
}

func (os orbitSet) Len() int {
	return len(os.set)
}
func (os orbitSet) Swap(i, j int) {
	os.set[i], os.set[j] = os.set[j], os.set[i]
}

func (os orbitSet) Less(i, j int) bool {
	// true if i < j
	a := os.set[i]
	b := os.set[j]

	vecA := pos{X: a.X - os.origin.X, Y: a.Y - os.origin.Y}
	vecB := pos{X: b.X - os.origin.X, Y: b.Y - os.origin.Y}

	angleA := angleWithTwelveOclock(vecA)
	angleB := angleWithTwelveOclock(vecB)

	if angleA == angleB {
		magA := math.Sqrt(float64(a.X + a.Y))
		magB := math.Sqrt(float64(b.X + b.Y))
		return magA > magB
	}
	return angleA < angleB

}

func angleWithTwelveOclock(p pos) float64 {
	α := float64(0)

	α = ((math.Pi / 2) + math.Atan2(float64(p.Y), float64(p.X)))

	// make sure are values are represented as clockwise from 12 oclock
	if α < 0 {
		α += math.Pi * 2
	}

	return α
}
