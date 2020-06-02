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
	"sync"
	"time"

	"./HAL"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments specified")
		return
	}
	debugModus := flag.Bool("debug", false, "enables debug output")
	confPath := flag.String("configuration", "", "The HAL-Programm file")
	flag.Parse()

	// Load the configuration JSON into a structure
	confStructure := readConfiguration(confPath)

	// Create the channels
	var listOfChannels []createdChannels
	for _, connection := range confStructure.HALVerbindungen {
		tmpChannel := make(chan float64, 4)
		tmpCreatedChannels := createdChannels{connection.From.Prozessor, connection.From.Channel, connection.To.Prozessor, connection.To.Channel, tmpChannel}
		listOfChannels = append(listOfChannels, tmpCreatedChannels)
	}

	// Configure and Start HAL
	var wg sync.WaitGroup
	wg.Add(4)
	for _, prozessor := range confStructure.HALProzessoren {
		instructions := readFile(&prozessor.Directory)
		var prozessorConnections []HAL.Connection
		for _, connection := range listOfChannels {
			if prozessor.Prozessor == connection.fromProzessor {
				prozessorConnections = append(prozessorConnections, HAL.Connection{connection.fromPort, connection.Channel, "from"})
			}
			if prozessor.Prozessor == connection.toProzessor {
				prozessorConnections = append(prozessorConnections, HAL.Connection{connection.toPort, connection.Channel, "to"})
			}
		}
		//fmt.Println("Prozessor: ", prozessor.Prozessor, "\nConnections:\n", prozessorConnections, "\n")
		go HAL.HalStart(instructions, *debugModus, prozessor.Prozessor, prozessorConnections, &wg)
		time.Sleep(0 * time.Second)

	}
	wg.Wait()

	fmt.Println("Main: Completed")
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

type createdChannels struct {
	fromProzessor int
	fromPort      int
	toProzessor   int
	toPort        int
	Channel       chan float64
}

func readConfiguration(input *string) jsonConfiguration {
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
