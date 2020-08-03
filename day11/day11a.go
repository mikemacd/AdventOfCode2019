package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type programType struct {
	instructions map[int]int
	relativeBase int
	pc           int // current program counter / instruction pointer
	hull         hullType
}

type positionType struct {
	X int
	Y int
}

type positionTypeSet []positionType

type orbitSet struct {
	origin positionType
	set    positionTypeSet
}

type hullType struct {
	min       positionType
	max       positionType
	robot     robotType
	positions map[positionType]int
}

type robotType struct {
	pos     positionType
	bearing positionType
}

var debug = false

var itteration = 0
var done = make(chan int)

var (
	bearU = positionType{X: 0, Y: 1}
	bearL = positionType{X: -1, Y: 0}
	bearD = positionType{X: 0, Y: -1}
	bearR = positionType{X: 1, Y: 0}
)

func main() {
	program := programType{
		instructions: map[int]int{},
		pc:           0,
		relativeBase: 0,
		hull: hullType{
			min: positionType{0, 0},
			max: positionType{0, 0},
			robot: robotType{
				pos:     positionType{X: 0, Y: 0},
				bearing: bearU,
			},
			positions: map[positionType]int{
				{X: 0, Y: 0}: 0,
			},
		},
	}
	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}

	if len(os.Args) < 3 {
		log.Fatalf("Missing parameter, provide input value!")
		return
	}

	program.readProgram(os.Args[1])

	var pipeSize int = 10000

	// each 'computer' will have an input
	aIn := make(chan int, pipeSize)
	aOut := make(chan int, pipeSize)

	aIn <- func() int { a, _ := strconv.Atoi(os.Args[2]); return a }() // run the test program

	// program.print()

	wg := sync.WaitGroup{}
	wg.Add(1)

	//	program.hull.print()

	go func() {
		program.run(0, aIn, aOut)
		done <- 1
		return
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("R:%v", r)
			fmt.Printf("size of positions: %d \n", len(program.hull.positions))
			for v := range aOut {
				fmt.Printf("LEFT OVER:%v\n", v)
			}
		}
	}()

	// inter program pipe
	go func(program *programType) {
		for {
			select {
			case colour := <-aOut:
				
				// do turtle stuff here
				// fmt.Printf("POSITIONS:%v\n",program.hull.positions)
				// fmt.Printf("BEFORE pos:%v BEAR:%v\n", program.hull.robot.pos, program.hull.robot.bearing)

				// paint
				program.hull.positions[program.hull.robot.pos] = colour

				direction := <-aOut
				if direction == 99 {
					done <- 1
					return
				}

				// turn
				program.hull.robot.bearing = directionTransformation[direction][program.hull.robot.bearing]

				newPos := positionType{
					X: program.hull.robot.pos.X + program.hull.robot.bearing.X,
					Y: program.hull.robot.pos.Y + program.hull.robot.bearing.Y,
				}
				// fmt.Printf("NEW pos:%v\n", newPos)

				// step
				program.hull.robot.pos = newPos

				// have we painted here before? If not, create a panel.
				if _, ok := program.hull.positions[program.hull.robot.pos]; !ok {
					program.hull.positions[program.hull.robot.pos] = 0
				}

				// keep track of the bounds for printing
				if program.hull.robot.pos.X > program.hull.max.X {
					program.hull.max.X = program.hull.robot.pos.X
				}
				if program.hull.robot.pos.Y > program.hull.max.Y {
					program.hull.max.Y = program.hull.robot.pos.Y
				}
				if program.hull.robot.pos.X < program.hull.min.X {
					program.hull.min.X = program.hull.robot.pos.X
				}
				if program.hull.robot.pos.Y < program.hull.min.Y {
					program.hull.min.Y = program.hull.robot.pos.Y
				}

				// fmt.Printf("AFTER robot:%v\n", program.hull.robot)

				// fmt.Fprintf(os.Stderr, "positionType:%v is colour:%v\n%v\n", program.hull.robot.positionType, program.hull.positionTypeitions[program.hull.robot.positionType],program.hull.positionTypeitions )

				aIn <- program.hull.positions[program.hull.robot.pos]

				// fmt.Printf("\nitteration:%d\n", itteration)
				// program.hull.print()

				itteration++

			case <-done:
				wg.Done()
				log.Printf("DONE\n\n")
				return
			default:
				// nop
			}
		}
	}(&program)

	wg.Wait()

	fmt.Printf("panels:%v\n", program.hull.positions)
	fmt.Printf("panels seen: %d\n", len(program.hull.positions))
	program.hull.print()

}

// positionTypeition is the positionTypeition to read the parameter from
func (p *programType) memGet(positionTypeition, mode int) int {
	positionTypeition = p.instructions[positionTypeition]

	rv := 0
	switch mode {
	case 0:
		rv = p.instructions[positionTypeition]
	case 1:
		rv = positionTypeition
	case 2:
		rv = p.instructions[positionTypeition+p.relativeBase]
	default:
		log.Printf("BOGUS MODE:%d\n", mode)
	}

	if debug {
		fmt.Printf("MEMget  positionType:%d mode:%d rv:%d RB:%d\n", positionTypeition, mode, rv, p.relativeBase)
	}
	return rv
}

func (p *programType) memSet(positionTypeition, mode, value int) {
	// mode 0: positionTypeition
	positionType := p.instructions[positionTypeition]

	// mode 1: positionTypeition
	if mode == 1 {
		positionType = positionTypeition
		log.Fatalf("BAD MODE (1) FOR MEMSET\n")
	}

	// mode 2: relative
	if mode == 2 {
		positionType = positionType + p.relativeBase
	}

	p.instructions[positionType] = value
	if debug {
		fmt.Printf("MEMSET  positionType:%d mode:%d value:%d RB:%d\n", positionType, mode, value, p.relativeBase)
	}
}

// opcode 1
func (p *programType) add(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	modeC := (modes / 100) % 10
	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)

	result := operandA + operandB

	p.memSet(p.pc+3, modeC, result)
	p.pc += 4

	if debug {
		fmt.Printf("id:%d\t\t\tadd:%v + %v = %v\n", id, operandA, operandB, result)
	}
}

// opcode 2
func (p *programType) mul(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	modeC := (modes / 100) % 10
	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)

	result := operandA * operandB

	p.memSet(p.pc+3, modeC, result)
	p.pc += 4

	if debug {
		fmt.Printf("id:%d\t\t\tmul:%v * %v = %v\n", id, operandA, operandB, result)
	}
}

// opcode 3
func (p *programType) input(id int, modes int, input chan int) {
	modeA := (modes / 1) % 10

	msg := <-input

	p.memSet(p.pc+1, modeA, msg)

	p.pc += 2

	if debug {
		fmt.Printf("id:%d\tinput\t\tmodes:%d modeA: %d msg:%d positionType:%d \n", id, modes, modeA, msg, p.pc+1)
	}

}

// opcode 4
func (p *programType) output(id int, modes int, output chan int) {
	modeA := (modes / 1) % 10
	parameter := p.memGet(p.pc+1, modeA)

	output <- parameter

	p.pc += 2

	if debug {
		fmt.Printf("id:%d\t\t\tmodes:%d modeA: %d output:%v\n", id, modes, modeA, parameter)
	}
}

// opcode 5
func (p *programType) jumpIfTrue(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	operandA := p.memGet(p.pc+1, modeA)
	positionTypeition := p.memGet(p.pc+2, modeB)

	if operandA != 0 {
		p.pc = positionTypeition
	} else {
		p.pc += 3
	}
	if debug {
		fmt.Printf("id:%d\t\t\tJiT:%v != 0 ? %v => positionType:%v\n", id, operandA, operandA != 0, positionTypeition)
	}
}

// opcode 6
func (p *programType) jumpIfFalse(id, modes int) {
	defer func() {
	}()
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	positionTypeition := p.memGet(p.pc+2, modeB)

	if operandA == 0 {
		p.pc = positionTypeition
	} else {
		p.pc += 3
	}

	if debug {
		fmt.Printf("id:%d\t\t\tJiF:%v == 0 ? %v => positionType:%v\n", id, operandA, operandA == 0, positionTypeition)
	}
}

// opcode 7
func (p *programType) lessThan(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	modeC := (modes / 100) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)

	rv := 0
	if operandA < operandB {
		rv = 1
	}

	p.memSet(p.pc+3, modeC, rv)
	p.pc += 4
	if debug {
		fmt.Printf("id:%d\t\t\tLT:%v < %v ? %v => positionType:%v\n", id, operandA, operandB, operandA < operandB, p.pc+3)
	}
}

// opcode 8
func (p *programType) equalTo(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	modeC := (modes / 100) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)

	rv := 0
	if operandA == operandB {
		rv = 1
	}

	p.memSet(p.pc+3, modeC, rv)
	p.pc += 4

	if debug {
		fmt.Printf("id:%d\t\t\tLT:%v == %v ? %v => positionType:%v\n", id, operandA, operandB, operandA == operandB, p.pc+3)
	}
}

// opcode 9
func (p *programType) adjustRelativeBase(id, modes int) {
	originalRB := p.relativeBase /// just for debug
	modeA := (modes / 1) % 10
	adjustment := p.memGet(p.pc+1, modeA)

	p.relativeBase += adjustment
	p.pc += 2
	if debug {
		fmt.Printf("ADJ_RB modeA: %d Orb:%d adjustment:%d new relativebase: %d\n", modeA, originalRB, adjustment, p.relativeBase)
	}
}

func (p *programType) readProgram(file string) {
	// raw reading of the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can't read file: %v\n", file)
		panic(err)
	}

	// take the read file and convert it from strings to ints
	for i, num := range bytes.Split([]byte(strings.TrimSpace(string(data))), []byte(",")) {
		code, err := strconv.Atoi(string(num))
		if err != nil {
			log.Fatalf("Could not convert opcode %v to integer. %v\n", num, err)
		}

		p.instructions[i] = code
	}

	// index of current opcode
	p.pc = 0
	p.relativeBase = 0
}

func (p *programType) run(id int, in chan int, out chan int) {
	opcode := 0
	for opcode != 99 {
		if debug {
			p.print()
		}

		instruction := p.instructions[p.pc]
		modes := instruction / 100
		opcode = instruction % 100

		ocLabel := []string{"nop", "add", "mul", "input", "output", "jumpIfTrue", "jumpIfFalse", "lessThan", "equalTo", "adjRB", "exit"}
		for i := 0; i < 100; i++ {
			ocLabel = append(ocLabel, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		}

		if debug {
			fmt.Fprintf(os.Stdout, "id:%d pc:%d oc: %d %v\n", id, p.pc, opcode, ocLabel[opcode])
		}

		switch opcode {
		case 1:
			p.add(id, modes)
		case 2:
			p.mul(id, modes)
		case 3:
			p.input(id, modes, in)
		case 4:
			p.output(id, modes, out)
		case 5:
			p.jumpIfTrue(id, modes)
		case 6:
			p.jumpIfFalse(id, modes)
		case 7:
			p.lessThan(id, modes)
		case 8:
			p.equalTo(id, modes)
		case 9:
			p.adjustRelativeBase(id, modes)
		case 99:
			if debug {
				spew.Fdump(os.Stdout, p)
			}
			done <- 1

			// outpt diag code
			return
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}

}

func (p *programType) print() {
	text := ""
	// To store the keys in slice in sorted order
	var keys []int
	for k := range p.instructions {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// To perform the opertion you want
	for _, k := range keys {
		text += fmt.Sprintf("[%d]:%v\t", k, p.instructions[k])
	}
	log.Printf("%s\n", text)
}

var directionTransformation = map[int]map[positionType]positionType{
	// left
	0: {
		bearU: bearL,
		bearL: bearD,
		bearD: bearR,
		bearR: bearU,
	},
	1: {
		bearU: bearR,
		bearR: bearD,
		bearD: bearL,
		bearL: bearU,
	},
}

func (h hullType) print() {
	for j := h.max.Y + 1; j >= h.min.Y-1; j-- {
		for i := h.min.X - 1; i <= h.max.X+1; i++ {
			output := " "
			if h.positions[positionType{i, j}] == 1 {
				output = "#"
			}
			if v,ok:=h.positions[positionType{i, j}];ok&& v == 0 {
				output = "."
			}
			if (h.robot.pos == positionType{X: i, Y: j}) {
				switch h.robot.bearing {
				case bearU:
					output = "^"
				case bearR:
					output = ">"
				case bearD:
					output = "v"
				case bearL:
					output = "<"
				}
			}
			fmt.Printf("%s", output)
		}
		fmt.Printf("\n")
	}
	//fmt.Printf("size of positions: %d \n",len(h.positions))
}
