package main

import (
	"strings"
	"bytes"
	"io/ioutil"
	"fmt"
 	"os"
	"strconv"
 )

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}

	sum := 0
	for i, line := range bytes.Split(data, []byte("\n")) {
		number := strings.TrimSpace(string(line))
		if number != "" {
			n, err := strconv.Atoi(string(number))
			if err != nil {
				fmt.Printf("Could not convert line %d (%s) to int\n", i, number)
			}
			sum += fuelForFuel(n)
		}
	}
	fmt.Printf("Sum: %d\n", sum)

}

func fuelForFuel(f int) (int) {
	f4f := (f/3)-2
	if f4f >0 {
		return f4f+fuelForFuel(f4f)
	}
	return 0
}
 