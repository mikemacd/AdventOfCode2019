package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	var sum int

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}

	sum = 0
	for i, line := range bytes.Split(data, []byte("\n")) {
		number := strings.TrimSpace(string(line))
		if number != "" {
			n, err := strconv.Atoi(string(number))
			if err != nil {
				fmt.Printf("Could not convert line %d (%s) to int\n", i, number)
			}
			sum += (n / 3) - 2
		}
	}

	fmt.Printf("Sum: %d\n", sum)

}
