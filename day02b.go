package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var debug = false

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	// read the input from the specified file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}

	// Process the input
	var input []int
 	for i, line := range bytes.Split(data, []byte(",")) {
		number := strings.TrimSpace(string(line))
		if number != "" {
			n, err := strconv.Atoi(string(number))
			if err != nil {
				fmt.Printf("Could not convert line %d (%s) to int\n", i, number)
			}
			input = append(input,n)
		}
	}
	if debug {
		fmt.Printf("INPUT: %+v\n", input)
	}

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			func() {
				defer func() {
					// recover from panic if one occured. Set err to nil otherwise.
					// /*
					if r := recover(); r != nil {
						fmt.Printf("NOPE: noun:%d verb:%d \n", noun, verb)
					}

					// */
				}()

				code := make([]int, len(input))
				copy(code, input)
  				code[1] = noun
				code[2] = verb


				for i := 0; i < len(code); i += 4 {
					opcode := code[i]
					if opcode == 99 {
						if debug {
							fmt.Printf("EXIT\n")
							fmt.Printf("noun:%d verb:%d i:%d\topcode:%d\toperand1:%d\toperand2:%d\tdest:%d\n",  noun,verb, i, opcode, -1, -1, code[i+3])
						}
						break
					}
					operand1 := code[code[i+1]]
					operand2 := code[code[i+2]]
					if debug {
						fmt.Printf("noun:%d verb:%d i:%d\topcode:%d\toperand1:%d\toperand2:%d\tdest:%d\n", noun,verb,  i, opcode, operand1, operand2, code[i+3])
					}
					switch opcode {
					case 1:
						if debug {
							fmt.Printf("noun:%d verb:%d %d = %d + %d = %d\n", noun,verb,  code[i+3], operand1, operand2, operand1+operand2)
						}

						code[code[i+3]] = operand1 + operand2

					case 2:
						if debug {
							fmt.Printf("noun:%d verb:%d %d = %d * %d = %d\n", noun,verb,  code[i+3], operand1, operand2, operand1-operand2)
						}

						code[code[i+3]] = operand1 * operand2
					default:
						fmt.Printf("DEFAULT: i:%d oc:%d", i, opcode)
					}

					if debug {
						fmt.Printf("noun:%d verb:%d %+v\n", noun, verb, code)
					}
				}
				fmt.Printf("noun:%d verb:%d Result: %d\n", noun, verb, code[0])

				if code[0] == 19690720 {
					fmt.Printf("YES! %+v\n", 100*noun+verb)
					os.Exit(0)
				}
			}()
		}
	}

}
