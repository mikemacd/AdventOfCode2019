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

	"github.com/davecgh/go-spew/spew"
)

type programType struct {
	instructions map[int]int
	relativeBase int
	pc           int // current program counter / instruction pointer
}

var debug = false

var program programType

func main() {
	program := programType{
		instructions: map[int]int{},
		pc:           0,
		relativeBase: 0,
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
	program.run(0, aIn, aOut)

	close(aOut)

	for msg := range aOut {
		fmt.Printf("%v\n", msg)
	}
}

// position is the position to read the parameter from
func (p *programType) memGet(position, mode int) int {
	position = p.instructions[position]

	rv := 0
	switch mode {
	case 0:
		rv = p.instructions[position]
	case 1:
		rv = position
	case 2:
		rv = p.instructions[position+p.relativeBase]
	default:
		log.Printf("BOGUS MODE:%d\n", mode)
	}

	if debug {
		fmt.Printf("MEMget  pos:%d mode:%d rv:%d RB:%d\n", position, mode, rv, p.relativeBase)
	}
	return rv
}

func (p *programType) memSet(position, mode, value int) {
	// mode 0: position
	pos := p.instructions[position]

	// mode 1: position
	if mode == 1 {
		pos = position
		log.Fatalf("BAD MODE (1) FOR MEMSET\n")
	}

	// mode 2: relative
	if mode == 2 {
		pos = pos + p.relativeBase
	}

	p.instructions[pos] = value
	if debug {
		fmt.Printf("MEMSET  pos:%d mode:%d value:%d RB:%d\n", pos, mode, value, p.relativeBase)
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
		fmt.Printf("id:%d\tinput\t\tmodes:%d modeA: %d msg:%d pos:%d \n", id, modes, modeA, msg, p.pc+1)
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
	position := p.memGet(p.pc+2, modeB)

	if operandA != 0 {
		p.pc = position
	} else {
		p.pc += 3
	}
	if debug {
		fmt.Printf("id:%d\t\t\tJiT:%v != 0 ? %v => pos:%v\n", id, operandA, operandA != 0, position)
	}
}

// opcode 6
func (p *programType) jumpIfFalse(id, modes int) {
	defer func() {
	}()
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	position := p.memGet(p.pc+2, modeB)

	if operandA == 0 {
		p.pc = position
	} else {
		p.pc += 3
	}

	if debug {
		fmt.Printf("id:%d\t\t\tJiF:%v == 0 ? %v => pos:%v\n", id, operandA, operandA == 0, position)
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
		fmt.Printf("id:%d\t\t\tLT:%v < %v ? %v => pos:%v\n", id, operandA, operandB, operandA < operandB, p.pc+3)
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
		fmt.Printf("id:%d\t\t\tLT:%v == %v ? %v => pos:%v\n", id, operandA, operandB, operandA == operandB, p.pc+3)
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
				spew.Fdump(os.Stdout, program)
			}
			// outpt diag code
			return
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}

}

// from
//		https://stackoverflow.com/questions/42028130/go-language-create-permutations
func permute(set []int) [][]int {
	permutations := [][]int{}
	if len(set) == 1 {
		return [][]int{set}
	}

	for i := range set {
		el := make([]int, len(set))
		copy(el, set)

		for _, perm := range permute(append(el[0:i], el[i+1:]...)) {
			permutations = append(permutations, append([]int{set[i]}, perm...))
		}
	}
	return permutations
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
