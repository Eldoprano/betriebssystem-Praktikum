package main

import(
	"fmt"
	"flag"
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
	"./HAL"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments specified")
		return
	}
	d := flag.Bool("debug", false, "enables debug output")
	input := flag.String("input", "", "The HAL-Programm file")
	flag.Parse()
	m := readFile(input)

	HAL.HalStart(m, *d)
}

func readFile(input *string) map[int]string {
	var m map[int]string
	m = make(map[int]string)
	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		lineNumber, err := strconv.Atoi(s[0])
		if err != nil {
			log.Fatal("Line did not start with a number")
		}
		m[lineNumber] = strings.Join(s[1:], " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m

}