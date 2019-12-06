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
	var code []int
 	for i, line := range bytes.Split(data, []byte(",")) {
		number := strings.TrimSpace(string(line))
		if number != "" {
			n, err := strconv.Atoi(string(number))
			if err != nil {
				fmt.Printf("Could not convert line %d (%s) to int\n", i, number)
			}
			code = append(code,n)
		}
	}
	if debug {
		fmt.Printf("%+v\n", code)
	}

 	code[1]=12
 	code[2]=2

	for i:=0; i< len(code); i+=4   {
		opcode := code[i]
		if opcode==99 {
			if debug {
				fmt.Printf("i:%d\topcode:%d\toperand1:%d\toperand2:%d\tdest:%d\n", i, opcode, -1, -1, code[i+3])
			}
			break
		}
		operand1 := code[code[i+1]]
		operand2 := code[code[i+2]]
		if debug {
			fmt.Printf("i:%d\topcode:%d\toperand1:%d\toperand2:%d\tdest:%d\n", i, opcode, operand1, operand2, code[i+3])
		}
		switch opcode {
		case 1:
			if debug {
				fmt.Printf("%d = %d + %d = %d\n", code[i+3], operand1, operand2, operand1+operand2)
			}
			code[code[i+3]] = operand1 + operand2

		case 2:
			if debug {
				fmt.Printf("%d = %d - %d = %d\n", code[i+3], operand1, operand2, operand1-operand2)
			}
			code[code[i+3]] = operand1 * operand2

		}

		if debug {
			fmt.Printf("%+v\n", code)
		}
	}

	fmt.Printf("%+v\n", code[0])

}
