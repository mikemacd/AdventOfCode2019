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

	if debug {
		fmt.Printf("%v \n", program)
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
				// add
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
				// multiply
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
		case 3:
			{
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

				// save
				operandA := program[programCounter+1]
				program[operandA] = input

				programCounter += 2
			}
		case 4:
			{
				// save
				operandA := program[programCounter+1]
				fmt.Printf("output: %d", program[operandA])
				programCounter += 2
			}
		case 99:
			break
		default:
			log.Fatalf("Unexpected opcode: %+v\n", opcode)
		}
	}

}
