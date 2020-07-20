package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/davecgh/go-spew/spew"
)

var debug = false

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
	var program []int
	for _, num := range bytes.Split([]byte(strings.TrimSpace(string(data))), []byte(",")) {
		code, err := strconv.Atoi(string(num))
		if err != nil {
			log.Fatalf("Could not convert opcode %v to integer. %v\n", num, err)
		}
		program = append(program, code)
	}

	// input modifications specified in problem
	program[1] = 12
	program[2] = 2

	if debug {
		fmt.Printf("%v \n", program[:30])
	}

	// index of current opcode
	programCounter := 0
	opcode := 0

	for opcode != 99 {
		if debug {
			fmt.Printf("PC:%+v\t", programCounter)
		}
		opcode = program[programCounter]

		switch opcode {
		case 1:
			{
				operandA := program[program[programCounter+1]]
				operandB := program[program[programCounter+2]]
				position := program[programCounter+3]
				result := operandA + operandB
				program[position] = result

				if debug {
					fmt.Printf("%v + %v = %v => %v \t %v \n", operandA, operandB, result, position, program[:30])
				}
				programCounter += 4
			}
		case 2:
			{
				operandA := program[program[programCounter+1]]
				operandB := program[program[programCounter+2]]
				position := program[programCounter+3]
				result := operandA * operandB
				program[position] = result

				if debug {
					fmt.Printf("%v * %v = %v => %v \t %v \n", operandA, operandB, result, position, program[:30])
				}
				programCounter += 4
			}
		case 99:
			break
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}

	fmt.Printf("\n%v\n\n", program[0])

}
