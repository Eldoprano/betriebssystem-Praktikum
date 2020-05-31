package main

import (
	"bufio"
	"encoding/json"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"./HAL"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments specified")
		return
	}
	d := flag.Bool("debug", false, "enables debug output")
	confPath := flag.String("configuration", "", "The HAL-Programm file")
	progPath := flag.String("input", "", "The HAL-Programm file")
	flag.Parse()

	// Load the configuration JSON into a structure
	confStructure := readConfiguration(confPath)
	fmt.Println(confStructure)

	// Load the Programm file into a Map
	m := readFile(progPath)

	// Start the HAL
	HAL.HalStart(m, *d)
}

type jsonConfiguration struct {
	HALProzessoren  []prozessor  `json:"HALProzessoren"`
	HALVerbindungen []verbindung `json:"HALVerbindungen"`
}

type prozessor struct {
	Prozessor int    `json:"prozessor"`
	Directory string `json:"directory"`
}

type verbindung struct {
	From direction
	To   direction
}

type direction struct {
	Direction string
	Prozessor int
	Channel   int
}

func readConfiguration(input *string) jsonConfiguration{
	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	// Here we convert the JSON into a structure
	var result jsonConfiguration
	json.Unmarshal([]byte(byteValue), &result)
	return result
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
