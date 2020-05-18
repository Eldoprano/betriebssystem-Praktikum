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

	interpret(debugModus, HALInstructions)

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
	var accumulator float64 = 0
	register := make(map[int]float64)
	inOut := make(map[int]float64)

	for programmCounter < len(instructions) {
		// Save actual instruction
		actualInstruction := strings.ToUpper(instructions[programmCounter][0])

		if programmCounter == 0 && actualInstruction != "START" {
			fmt.Println("Instruction set needs to beginn with 'START'")
			os.Exit(3)
		}

		// Show command in case of Debug Modus
		switch actualInstruction {
		case "START":
			if debugModus {
				println(instructions[programmCounter][0], ":")
			}

		case "STOP":
			if debugModus {
				println(instructions[programmCounter][0], ":")
			}
			break

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

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "IN":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = inOut[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "LOAD":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = register[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "LOADNUM":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempFloat := strToFloat(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = tempFloat

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "STORE":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to float
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			register[tempInt] = accumulator

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "JUMPNEG":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			if accumulator < 0 {
				if tempInt >= len(instructions) {
					fmt.Println(strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter, "is jumping outside the barriers")
					os.Exit(3)
				}
				programmCounter = tempInt - 1 // Minus 1, because the programmCounter gets incremented at the end
			}

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "JUMPPOS":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			if accumulator > 0 {
				if tempInt >= len(instructions) {
					fmt.Println(strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter, "is jumping outside the barriers")
					os.Exit(3)
				}
				programmCounter = tempInt - 1 // Minus 1, because the programmCounter gets incremented at the end
			}

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "JUMPNULL":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			if accumulator == 0 {
				if tempInt >= len(instructions) {
					fmt.Println(strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter, "is jumping outside the barriers")
					os.Exit(3)
				}
				programmCounter = tempInt - 1 // Minus 1, because the programmCounter gets incremented at the end
			}

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "JUMP":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			if tempInt >= len(instructions) {
				fmt.Println(strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter, "is jumping outside the barriers")
				os.Exit(3)
			}
			programmCounter = tempInt - 1 // Minus 1, because the programmCounter gets incremented at the end

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "ADD":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator += register[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "ADDNUM":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempFloat := strToFloat(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator += tempFloat

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "SUB":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator -= register[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "MUL":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = accumulator * register[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "DIV":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempInt := strToInt(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = accumulator / register[tempInt]

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "SUBNUM":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempFloat := strToFloat(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator -= tempFloat

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "MULNUM":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempFloat := strToFloat(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = accumulator * tempFloat

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}

		case "DIVNUM":
			// Check if we got the second value
			if len(instructions[programmCounter]) < 2 {
				fmt.Println("Missing value for instruction", strings.ToUpper(instructions[programmCounter][0]), "in line", programmCounter)
				os.Exit(3)
			}
			// Convert second value to int
			tempFloat := strToFloat(instructions[programmCounter][1])

			// Execute the instruction:
			accumulator = accumulator / tempFloat

			// Show command in case of Debug Modus
			if debugModus {
				println(instructions[programmCounter][0], instructions[programmCounter][1], ":")
			}
		}
		if debugModus {
			printStatus(programmCounter, accumulator, register, inOut)
		}
		programmCounter++
	}

	if !debugModus {
		printStatus(programmCounter-1, accumulator, register, inOut)
	}
	//printStatus(programmCounter, accumulator, register, inOut)

}

func printStatus(programmCounter int, accumulator float64, register map[int]float64, inOut map[int]float64) {
	fmt.Println(" ______________________________________________________")
	fmt.Printf("| Programm Counter: %-8d| Accumulator: %-12.4f|\n", programmCounter, accumulator)
	//fmt.Printf("|Programm Counter:", programmCounter)
	//fmt.Println("|Accumulator:", accumulator)
	for i := 0; i < len(inOut); i = i + 2 {
		fmt.Printf("| I/O %d: %-19.4f| I/O %d: %-18.4f|\n", i, inOut[i], i+1, inOut[i+1])
		//fmt.Println("|I/O", i, ":", inOut[i], "\t\tI/O", i+1, ":", inOut[i+1])
	}
	for i := 0; i < 16; i = i + 2 {
		fmt.Printf("| Register %2d: %-13.4f| Register %2d: %-12.4f|\n", i, register[i], i+1, register[i+1])
		//fmt.Println("|Register", i, ":", register[i], "\t\tRegister", i+1, ":", register[i+1])
	}
	fmt.Print(" ______________________________________________________\n\n")
}

func strToInt(str string) int {
	tempInt, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Value needs to be an INT")
		os.Exit(3)
	}
	return tempInt
}

func strToFloat(str string) float64 {
	tempFloat, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Value needs to be an FLOAT")
		os.Exit(3)
	}
	return tempFloat
}
