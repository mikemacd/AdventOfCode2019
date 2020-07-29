package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type programType struct {
	instructions []int
	pc           int // current program counter / instruction pointer
}

var debug = false

var program programType

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}

	program.readProgram(os.Args[1])

	phaseSeqs := permute([]int{9, 8, 7, 6, 5})
	// phaseSeqs := [][]int{{9, 8, 7, 6, 5}}
	//  phaseSeqs := [][]int{{9,7,8,5,6}}
	var wg sync.WaitGroup
	var pipeSize int = 1000
	{
		biggestSignal := 0
		for _, sequence := range phaseSeqs {
			// fmt.Printf("sequence: %v\n", sequence)
			var signal int
			var ampA, ampB, ampC, ampD, ampE programType

			done := make(chan int)

			// each 'computer' will have an input
			aIn := make(chan int, pipeSize)
			bIn := make(chan int, pipeSize)
			cIn := make(chan int, pipeSize)
			dIn := make(chan int, pipeSize)
			eIn := make(chan int, pipeSize)
			aOut := make(chan int, pipeSize)
			bOut := make(chan int, pipeSize)
			cOut := make(chan int, pipeSize)
			dOut := make(chan int, pipeSize)
			eOut := make(chan int, pipeSize)

			ampA = program
			ampB = program
			ampC = program
			ampD = program
			ampE = program
			ampA.instructions = make([]int, len(program.instructions))
			ampB.instructions = make([]int, len(program.instructions))
			ampC.instructions = make([]int, len(program.instructions))
			ampD.instructions = make([]int, len(program.instructions))
			ampE.instructions = make([]int, len(program.instructions))
			copy(ampA.instructions, program.instructions)
			copy(ampB.instructions, program.instructions)
			copy(ampC.instructions, program.instructions)
			copy(ampD.instructions, program.instructions)
			copy(ampE.instructions, program.instructions)

			programCopy := []programType{ampA, ampB, ampC, ampD, ampE}

			wg.Add(1)
			go func() {
				programCopy[0].run(0, aIn, aOut)
			}()
			go func() {
				programCopy[1].run(1, bIn, bOut)
			}()
			go func() {
				programCopy[2].run(2, cIn, cOut)
			}()
			go func() {
				programCopy[3].run(3, dIn, dOut)
			}()
			go func() {
				defer wg.Done()
				programCopy[4].run(4, eIn, eOut)
				done <- 1
			}()

			// inter program pipe
			go func() {
				for {
					select {
					case msg := <-aOut:
						// log.Printf("a:%v\n\n", msg)
						bIn <- msg
					case msg := <-bOut:
						// log.Printf("b:%v\n\n", msg)
						cIn <- msg
					case msg := <-cOut:
						// log.Printf("c:%v\n\n", msg)
						dIn <- msg
					case msg := <-dOut:
						// log.Printf("d:%v\n\n", msg)
						eIn <- msg
					case msg := <-eOut:
						// log.Printf("e:%v\n\n", msg)
						signal = msg
						aIn <- msg
					case <-done:
						// log.Printf("DONE\n\n")
						return
					default:
						// nop
					}
				}
			}()

			// prime the inputs with the signal value from the sequence. In reverse order o prevent the phase being preempted.
			eIn <- sequence[4]
			dIn <- sequence[3]
			cIn <- sequence[2]
			bIn <- sequence[1]
			aIn <- sequence[0]

			// 'start' the network by providing a zero to computer A
			aIn <- 0

			wg.Wait()

			if signal > biggestSignal {
				biggestSignal = signal
			}
		}

		fmt.Printf("Biggest signal: %v", biggestSignal)
	}
}

func (p *programType) memGet(position, mode int) int {
	if mode == 0 {
		return p.instructions[p.instructions[position]]
	}

	return p.instructions[position]
}

func (p *programType) memSet(position, value int) {
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

func (p *programType) add(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA + operandB

	// fmt.Printf("id:%d\t\t\tadd:%v + %v = %v\n", id, operandA, operandB, result)

	p.memSet(position, result)

	p.pc += 4
}

func (p *programType) mul(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1) // we're just looking up the offset to be used later in memSet thus making the memset be in immediate mode

	result := operandA * operandB

	//fmt.Printf("id:%d\t\t\tmul:%v * %v = %v\n", id, operandA, operandB, result)

	p.memSet(position, result)

	p.pc += 4
}

func (p *programType) input(id int, modes int, input chan int) {
	defer func() {
		p.pc += 2
	}()

	msg := <-input

	// fmt.Printf("id:%d\t\t\tinput:%v\n", id, msg)

	p.memSet(p.memGet(p.pc+1, 1), msg)
}

func (p *programType) output(id int, modes int, output chan int) {
	defer func() {
		p.pc += 2
	}()

	modeA := (modes / 1) % 10

	operandA := p.memGet(p.pc+1, modeA)

	//	fmt.Printf("id:%d\t\t\toutput:%v\n", id, operandA)

	output <- operandA
}

func (p *programType) jumpIfTrue(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	position := p.memGet(p.pc+2, modeB)

	// fmt.Printf("id:%d\t\t\tJiT:%v != 0 ? %v => pos:%v\n", id, operandA, operandA != 0 , position)

	if operandA != 0 {
		p.pc = position
		return
	}
	p.pc += 3
}

func (p *programType) jumpIfFalse(id, modes int) {
	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	position := p.memGet(p.pc+2, modeB)

	// fmt.Printf("id:%d\t\t\tJiF:%v == 0 ? %v => pos:%v\n", id, operandA, operandA == 0 , position)

	if operandA == 0 {
		p.pc = position
		return
	}
	p.pc += 3
}

func (p *programType) lessThan(id, modes int) {
	defer func() {
		p.pc += 4
	}()

	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1)

	// fmt.Printf("id:%d\t\t\tLT:%v < %v ? %v => pos:%v\n", id, operandA, operandB, operandA < operandB , position)

	if operandA < operandB {
		p.memSet(position, 1)
		return
	}
	p.memSet(position, 0)
}

func (p *programType) equalTo(id, modes int) {
	defer func() {
		p.pc += 4
	}()

	modeA := (modes / 1) % 10
	modeB := (modes / 10) % 10

	operandA := p.memGet(p.pc+1, modeA)
	operandB := p.memGet(p.pc+2, modeB)
	position := p.memGet(p.pc+3, 1)

	// fmt.Printf("id:%d\t\t\tLT:%v == %v ? %v => pos:%v\n", id, operandA, operandB, operandA == operandB , position)

	if operandA == operandB {
		p.memSet(position, 1)

		return
	}
	p.memSet(position, 0)

}

func (p *programType) readProgram(file string) {
	// raw reading of the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can't read file: %v\n", file)
		panic(err)
	}

	// take the read file and convert it from strings to ints
	for _, num := range bytes.Split([]byte(strings.TrimSpace(string(data))), []byte(",")) {
		code, err := strconv.Atoi(string(num))
		if err != nil {
			log.Fatalf("Could not convert opcode %v to integer. %v\n", num, err)
		}

		p.instructions = append(p.instructions, code)
	}

	// index of current opcode
	p.pc = 0
}

func (p *programType) run(id int, in chan int, out chan int) {
	opcode := 0
	for opcode != 99 {
		// fmt.Printf("id:%d prog:%v\n", id, p.instructions)
		if debug {
			fmt.Fprint(os.Stdout, "\n\n\n\n\n")
			fmt.Fprintln(os.Stdout, "##############")
			fmt.Fprintf(os.Stdout, "pc: %d \n", p.pc)
			fmt.Fprintln(os.Stdout, "##############")
			for i, v := range p.instructions {
				fmt.Fprintf(os.Stdout, "%d:%v\n", i, v)
			}
			fmt.Fprintln(os.Stdout, "##############")
		}

		instruction := p.instructions[p.pc]
		modes := instruction / 100
		opcode = instruction % 100

		ocLabel := []string{"nop", "add", "mul", "input", "output", "jumpIfTrue", "jumpIfFalse", "lessThan", "equalTo", "exit"}
		for i := 0; i < 100; i++ {
			ocLabel = append(ocLabel, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		}
		//		fmt.Fprintf(os.Stdout, "id:%d oc: %d %v\n", id, opcode, ocLabel[opcode])

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
