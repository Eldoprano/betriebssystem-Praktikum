package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv" // Convert Strings to Int
	"strings" // Split strings
)

var debugModus = false

func main() {
	// Check if we are in debug Modus
	debugModus = isInDebugModus()

	// Get a Map filled with an array of instructions (strings)
	// It also checks if the instructions beginn with a number, and are in order
	// HALInstructions[2] = "[STORE, 8]"
	// HALInstructions[2][0] = "STORE"
	HALInstructions, fileError := readFile("HAL_Instructions")
	if fileError != nil {
		fmt.Println(fileError)
		os.Exit(3)
	}

	fmt.Println(HALInstructions[0])

}

func isInDebugModus() bool {

	// Read Arguments to find if we need to switch to debugMode
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Println("HAL Interpreter in Normalmodus")
		return false
	} else if argsWithoutProg[0] == "--debug" || argsWithoutProg[0] == "-d" {
		fmt.Println("HAL Interpreter in Debugmodus")
		return true
	} else {
		fmt.Println("Unrecognized arguments")
		fmt.Println("HAL Interpreter in Normalmodus")
		return false
	}
}

func readFile(fileName string) (map[int][]string, error) {
	// Read file content
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	// Close file when function returns something
	defer file.Close()

	// Fill the Map with the commands on the file
	commands := make(map[int][]string)

	scanner := bufio.NewScanner(file)
	orderChecker := 0
	for scanner.Scan() {
		actualLine := scanner.Text()

		// Check if first part of the line is a number
		num, err := strconv.Atoi(strings.Split(actualLine, " ")[0])
		if err != nil {
			err = errors.New("HAL Instruction needs to beginn with a number")
			return nil, err
		}

		// Check if the order is correct
		if num != orderChecker {
			err = errors.New("The HAL Instructions are not in order")
			return nil, err
		}
		orderChecker++

		// Saves line on the Dictionary
		commands[num] = strings.Split(actualLine, " ")[1:]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}
