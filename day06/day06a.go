package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// the tree is an emergent property of the data once it has been loaded
type tree map[string][]string

// var debug = false

func main() {
	orbits := tree{}

	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}

	// raw reading of the file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Can't read file: %v\n", os.Args[1])
		panic(err)
	}

	// take the read file and convert it from strings to ints
	for _, line := range bytes.Split(data, []byte("\n")) {
		planets := strings.Split(strings.TrimSpace(string(line)), ")")
		if err != nil {
			log.Fatalf("Could not extract planets from line. %v\n", err)
		}

		orbits[planets[0]] = append(orbits[planets[0]], planets[1])
	}

	spew.Dump("Result:")
	//	spew.Dump(orbits)

	//	fmt.Printf("nodes:%d\n", orbits.height("COM"))
	//	fmt.Printf("count:%d\n", orbits.count("COM", 0))
	fmt.Printf("p:%v\n", orbits.parent("B"))
	fmt.Printf("p:%v\n", orbits.parent("COM"))
	fmt.Printf("q:%v\n", orbits.ancestry("I"))
	fmt.Printf("q:%v\n", orbits.ancestry("L"))

}

func (t tree) height(node string) int {
	height := 1
	for _, child := range t[node] {
		height += t.height(child)
	}
	return height
	// height of COM is 1 + the heiht of all of the childrens children.
}

func (t tree) count(node string, height int) int {
	subcount := height

	for _, child := range t[node] {
		subcount += t.count(child, height+1)
	}

	return subcount
	// height of COM is 1 + the heiht of all of the childrens children.
}

func (t tree) parent(node string) string {
	for i, v := range t {
		for _, q := range v {
			if node == q {
				return i
			}
		}
	}
	return ""
}

func (t tree) ancestry(node string) []string {
	p := t.parent(node)

	if p == "" {
		return []string{node}
	}

	a := t.ancestry(p)

	return append(a, node)
}
