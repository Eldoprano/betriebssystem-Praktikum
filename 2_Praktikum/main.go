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

	// Get a String array filled with the Instructions
	// HALInstructions[0] = "00 Start"
	HALInstructions, fileError := readFile("HAL_Instructions")
	if fileError != nil {
		fmt.Println(fileError)
	} else {
		fmt.Println(HALInstructions[0])
	}
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

func readFile(fileName string) ([]string, error) {
	// Read file content
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	// Close file when function returns something
	defer file.Close()

	// Fill the String array with the commands on the file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		actualLine := scanner.Text()

		// Check if first part of the line is a number
		_, err := strconv.Atoi(strings.Split(actualLine, " ")[0])
		if err != nil {
			//fmt.Println(err)
			err = errors.New("HAL Instruction needs to beginn with a number")
			return nil, err
		}

		lines = append(lines, actualLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
