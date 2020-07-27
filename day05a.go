package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type programType struct {
	instructions []int
	pc           int // current program counter / instruction pointer
}

var program programType

func main() {

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
	for _, num := range bytes.Split([]byte(strings.TrimSpace(string(data))), []byte(",")) {
		code, err := strconv.Atoi(string(num))
		if err != nil {
			log.Fatalf("Could not convert opcode %v to integer. %v\n", num, err)
		}

		program.instructions = append(program.instructions, code)
	}

	// index of current opcode
	program.pc = 0
	opcode := 0

	for opcode != 99 {
		instruction := program.instructions[program.pc]
		modes := instruction / 100
		opcode = instruction % 100

		switch opcode {
		case 1:
			program.add(modes)
		case 2:
			program.mul(modes)
		case 3:
			program.input(modes)
		case 4:
			program.output(modes)
		case 99:
			// spew.Dump("EXIT:", program)

			os.Exit(0)
			break
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}

}

func (p *programType) memGet(position, mode int) int {
	if mode == 0 {
		return p.instructions[p.instructions[position]]
	}

	return p.instructions[position]
}

func (p *programType) memSet(position, value int) {
	p.instructions[position] = value
}

func (p *programType) add(modes int) {
	modeA := modes % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA + operandB

	p.memSet(position, result)

	p.pc += 4

}

func (p *programType) mul(modes int) {
	modeA := modes % 10
	modeB := (modes / 10) % 10
	//fmt.Printf("modes a:%d b:%d \n",modeA,modeB )

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA * operandB

	p.memSet(position, result)

	p.pc += 4
}

func (p *programType) input(modes int) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("while reading input: %v", err)
	}

	input, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		log.Fatalf("Bad number: %v -- %v", text, err)
	}

	p.memSet(p.memGet(p.pc+1, 1), input)

	p.pc += 2
}

func (p *programType) output(modes int) {
	modeA := modes % 10

	operandA := p.memGet(p.pc+1, modeA)
	fmt.Printf("output: %d\n", operandA)
	p.pc += 2
}
