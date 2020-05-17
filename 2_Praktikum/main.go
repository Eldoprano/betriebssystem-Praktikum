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

	interpret(true, HALInstructions)

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

// ___________________HAL INTERPRETER____________________

func interpret(debugModus bool, instructions map[int][]string) {
	// Variables
	var programmCounter int = 0
	var accumulator float32 = 0
	register := make(map[int]float32)
	inOut := make(map[int]float32)

	register[1] = 2

	// Create map
	instrSET := make(map[string]string)
	instrSET["START"] = "N"
	instrSET["STOP"] = "N"
	instrSET["OUT"] = "s"
	instrSET["IN"] = "s"
	instrSET["LOAD"] = "r"
	instrSET["LOADNUM"] = "k"
	instrSET["STORE"] = "r"
	instrSET["JUMPNEG"] = "a"
	instrSET["JUMPPOS"] = "a"
	instrSET["JUMPNULL"] = "a"
	instrSET["JUMP"] = "a"
	instrSET["ADD"] = "r"
	instrSET["ADDNUM"] = "k"
	instrSET["SUB"] = "r"
	instrSET["MUL"] = "a"
	instrSET["DIV"] = "a"
	instrSET["SUBNUM"] = "a"
	instrSET["MULNUM"] = "a"
	instrSET["DIVNUM"] = "a"

	for programmCounter < len(instructions) {
		actualInstruction := strings.ToUpper(instructions[programmCounter][0])

		if instrSET[actualInstruction] == "" {
			fmt.Println("Instruction", strings.ToUpper(instructions[programmCounter][0]), "not found")
			os.Exit(3)
		}
		if programmCounter == 0 && actualInstruction != "START" {
			fmt.Println("Instruction set needs to beginn with 'START'")
			os.Exit(3)
		}
		switch actualInstruction {
		case "START":
			if debugModus {
				fmt.Println("START-desu!")
			}
		case "STOP":
			if debugModus {
				fmt.Println("STOP-desu!")
			}
		case "OUT":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			inOut[tempInt] = accumulator
			if debugModus {
				fmt.Println("OUT-desu!")
			}

		case "IN":
			if debugModus {
				fmt.Println("IN-desu!")
			}
		case "LOAD":
			if debugModus {
				fmt.Println("LOAD-desu!")
			}
		case "LOADNUM":
			if debugModus {
				fmt.Println("LOADNUM-desu!")
			}
		case "STORE":
			if debugModus {
				fmt.Println("STORE-desu!")
			}
		case "JUMPNEG":
			if debugModus {
				fmt.Println("JUMPNEG-desu!")
			}
		case "JUMPPOS":
			if debugModus {
				fmt.Println("JUMPPOS-desu!")
			}
		case "JUMPNULL":
			if debugModus {
				fmt.Println("JUMPNULL-desu!")
			}
		case "JUMP":
			if debugModus {
				fmt.Println("JUMP-desu!")
			}
		case "ADD":
			if debugModus {
				fmt.Println("ADD-desu!")
			}
		case "ADDNUM":
			if debugModus {
				fmt.Println("ADDNUM-desu!")
			}
		case "SUB":
			if debugModus {
				fmt.Println("SUB-desu!")
			}
		case "MUL":
			if debugModus {
				fmt.Println("MUL-desu!")
			}
		case "DIV":
			if debugModus {
				fmt.Println("DIV-desu!")
			}
		case "SUBNUM":
			if debugModus {
				fmt.Println("SUBNUM-desu!")
			}
		case "MULNUM":
			if debugModus {
				fmt.Println("MULNUM-desu!")
			}
		case "DIVNUM":
			if debugModus {
				fmt.Println("DIVNUM-desu!")
			}
		}

		programmCounter++
	}
}

func strToInt(str string) int {
	tempInt, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Value needs to be an INT")
		os.Exit(3)
	}
	return tempInt
}
