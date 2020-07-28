package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type programType struct {
	instructions []int
	pc           int // current program counter / instruction pointer
}

var debug = false

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing parameter, provide file name!")
		return
	}

	program := readProgram(os.Args[1])
	_ = program
	//	spew.Dump(permute([]int{0, 1, 2}))

	phaseSeqs:=permute([]int{0,1,2,3,4})
	_=phaseSeqs
	{
		biggestSignal := 0

		// phaseSeqs := [][]int{
		// 	{3,1,2,4,0},
		// 	 {4,3,2,1,0},
		// 	{0,1,2,3,4,},
		//  }
	
		for _,v := range phaseSeqs {
	//v := []int{3,1,2,4,0}
	//v := []int{4,3,2,1,0}
	//v := []int{0,1,2,3,4,}
	//v := []int{1,0,4,3,2,}

		signal := 0
		for _, phaseSequence := range v {
				// fmt.Printf("\nPS: %d\n", phaseSequence )
				// todo: make a copy of program
				programCopy:=program
				out := ghostRun(programCopy, []byte(fmt.Sprintf("%d\n%d\n",phaseSequence,signal)))
				newsignal,err := strconv.Atoi(out)
				if err!=nil {
					log.Fatalf("signal not understood: %v", err)
				}
				// fmt.Printf("\tsignal: %v \n", out)
				signal=newsignal
			}
			//fmt.Printf("\tsig: %v \n", signal)
			if signal>biggestSignal{
				biggestSignal=signal
			}
		}


		// phase, input signal
		// out := ghostRun(program, []byte("1\n0\n"))
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

func (p *programType) add(modes int) {
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

func (p *programType) mul(modes int) {
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

func (p *programType) input(modes int) {
	defer func() {
		p.pc += 2
	}()
	if debug {
		fmt.Println("input")
	}

	var text string
	_, err := fmt.Fscanf(os.Stdin, "%s", &text)
	if err != nil {
		log.Fatalf("while scaning for input: %v", err)
	}

	input, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		log.Fatalf("Bad number: %v -- %v", text, err)
	}

	p.memSet(p.memGet(p.pc+1, 1), input)
}

func (p *programType) output(modes int) {
	defer func() {
		p.pc += 2
	}()
	if debug {
		fmt.Println("output")
	}
	modeA := (modes / 1) % 10

	operandA := p.memGet(p.pc+1, modeA)

	fmt.Fprintf(os.Stdout, "%d", operandA)
}

func (p *programType) jumpIfTrue(modes int) {
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

func (p *programType) jumpIfFalse(modes int) {
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

func (p *programType) lessThan(modes int) {
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

func (p *programType) equalTo(modes int) {
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

func readProgram(file string) programType {
	var program programType

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

		program.instructions = append(program.instructions, code)
	}

	// index of current opcode
	program.pc = 0

	return program
}

func run(program programType) {
	opcode := 0
	for opcode != 99 {
		if debug {
			fmt.Fprint(os.Stdout, "\n\n\n\n\n")
			fmt.Fprintln(os.Stdout, "##############")
			fmt.Fprintf(os.Stdout, "pc: %d \n", program.pc)
			fmt.Fprintln(os.Stdout, "##############")
			for i, v := range program.instructions {
				fmt.Fprintf(os.Stdout, "%d:%v\n", i, v)
			}
			fmt.Fprintln(os.Stdout, "##############")
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

// ghostrun runs the program while providing input from the supplied parameter, and returning any output
func ghostRun(program programType, input []byte) string {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Printf("A\n")
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(input); err != nil {
		log.Printf("B\n")
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Printf("B\n")
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpfile

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	run(program)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = oldStdout // restoring the real stdout
	output := <-outC

	return output
}
