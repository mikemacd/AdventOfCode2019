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

	"github.com/davecgh/go-spew/spew"
)

type programTypeB struct {
	instructions []int
	pc           int // current program counter / instruction pointer
}

var debug = false

func main() {
	var program programTypeB

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
		if debug {
			fmt.Print("\n\n\n\n\n")
			fmt.Println("##############")
			fmt.Printf("pc: %d \n", program.pc)
			fmt.Println("##############")
			for i, v := range program.instructions {
				fmt.Printf("%d:%v\n", i, v)
			}
			fmt.Println("##############")
		}

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
		case 5:
			program.jumpIfTrue(modes)
		case 6:
			program.jumpIfFalse(modes)
		case 7:
			program.lessThan(modes)
		case 8:
			program.equalTo(modes)
		case 99:
			if debug {
				spew.Dump(program)
			}
			// outpt diag code
			os.Exit(0)
			break
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}
}

func (p *programTypeB) memGet(position, mode int) int {
	if mode == 0 {
		return p.instructions[p.instructions[position]]
	}

	return p.instructions[position]
}

func (p *programTypeB) memSet(position, value int) {
	defer func() {
		if r := recover(); r != nil {

			fmt.Println("##############")
			fmt.Println("Recovered in f", r)
			fmt.Println("##############")

			for i, v := range p.instructions {
				fmt.Printf("%d:%v\n", i, v)
			}
			fmt.Println("##############")
			fmt.Printf("p:%d v:%d", position, value)
			os.Exit(1)
		}
	}()
	if len(p.instructions) <= position {
		for i := len(p.instructions) - 1; i <= position; i++ {
			p.instructions = append(p.instructions, -1)
		}
	}
	p.instructions[position] = value
}

func (p *programTypeB) add(modes int) {
	if debug {
		fmt.Println("add")
	}

	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA + operandB

	p.memSet(position, result)

	p.pc += 4
}

func (p *programTypeB) mul(modes int) {
	if debug {
		fmt.Println("mul")
	}
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10
	//fmt.Printf("modes a:%d b:%d \n",modeA,modeB )

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA * operandB

	p.memSet(position, result)

	p.pc += 4
}

func (p *programTypeB) input(modes int) {
	if debug {
		fmt.Println("input")
	}

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

func (p *programTypeB) output(modes int) {
	defer func() {
		p.pc += 2
	}()
	if debug {
		fmt.Println("output")
	}
	modeA := (modes / 1) % 10

	operandA := p.memGet(p.pc+1, modeA)
	fmt.Printf("output: %d\n", operandA)
}

func (p *programTypeB) jumpIfTrue(modes int) {
	if debug {
		fmt.Println("jumpIfTrue")
	}
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	position := p.memGet(p.pc+2, modeB)

	if debug {
		fmt.Printf("T A:%d != 0 ? %v @ %d", operandA, operandA != 0, position)
	}
	if operandA != 0 {
		p.pc = position
		return
	}
	p.pc += 3
}

func (p *programTypeB) jumpIfFalse(modes int) {
	if debug {
		fmt.Println("jumpIfFalse")
	}
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	position := p.memGet(p.pc+2, modeB)

	if debug {
		fmt.Printf("F A:%d == 0 ? %v @ %d", operandA, operandA == 0, position)
	}
	if operandA == 0 {
		p.pc = position
		return
	}
	p.pc += 3
}

func (p *programTypeB) lessThan(modes int) {
	defer func() {
		p.pc += 4
	}()
	if debug {
		fmt.Println("lessThan")
	}
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1)
	if debug {
		fmt.Printf("\t\t%d < %d = %b in pos %d\n", operandA, operandB, operandA < operandB, position)
	}
	if operandA < operandB {
		p.memSet(position, 1)
		return
	}
	p.memSet(position, 0)
}

func (p *programTypeB) equalTo(modes int) {
	defer func() {
		p.pc += 4
	}()

	if debug {
		fmt.Println("equalTo")
	}
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1)

	if operandA == operandB {
		p.memSet(position, 1)

		return
	}
	p.memSet(position, 0)

}
