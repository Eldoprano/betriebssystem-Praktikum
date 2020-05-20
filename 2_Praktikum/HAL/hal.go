package HAL

import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"bufio"
	"os"
	"time"
)

var (
	pc int = 0
	regList []register
	ioList []inAndOut
	a accu
	debug = false
)

type register struct {
	value  float64
	number int
}

func (r *register) setValue(newVal float64) {
	r.value = newVal
}

type inAndOut struct {
	value  float64
	number int
}

func (i *inAndOut) setValue(newVal float64) {
	i.value = newVal
}

type accu struct {
	value float64
}

func (a *accu) setValue(newVal float64) {
	a.value = newVal
}

//program counter

func HalStart(instr map[int]string, d bool) {
	if d {
		debug = true
		fmt.Println("Running in Debug Mode")
	}
	//initalize all stuff
	regList = InitRegisters()
	ioList = InitInAndOut()
	a = accu{value: 0}

	//parse instruction and parameters
	for i := 1; i <= len(instr); i++ {
		if debug {
			time.Sleep(2 * time.Second)
			fmt.Println("-----------------------------------")
		}
		value := instr[i]
		s := strings.Fields(value)
		parameter := strings.Join(s[1:], " ")
		//fmt.Println(parameter)
		switch instruction := s[0]; instruction {
		case "START":
			//start(parameter)
		case "STOP":
			stop()
		case "OUT":
			out(parameter)
		case "IN":
			infunc(parameter)
		case "LOAD":
			load(parameter)
		case "LOADNUM":
			loadnum(parameter)
		case "STORE":
			store(parameter)
		case "JUMPNEG":
			address, err := jumpneg(parameter)
			if err != nil {
				fmt.Println("not jumping")
			} else {
				i = address - 1
			}
		case "JUMPPOS":
			address, err := jumppos(parameter)
			if err != nil {
				fmt.Println("not jumping")
			} else {
				i = address - 1
			}
		case "JUMPNULL":
			address, err := jumpnull(parameter)
			if err != nil {
				//fmt.Println("not jumping")
			} else {
				i = address - 1
			}
		case "JUMP":
			address := jump(parameter)
			i = address - 1
		case "ADD":
			add(parameter)
		case "ADDNUM":
			addnum(parameter)
		case "SUB":
			sub(parameter)
		case "MUL":
			mul(parameter)
		case "DIV":
			div(parameter)
		case "SUBNUM":
			subnum(parameter)
		case "MULNUM":
			mulnum(parameter)
		case "DIVNUM":
			divnum(parameter)
		default:
			fmt.Println("Invalid Instruction", instruction)
			os.Exit(2)
		}

		
	}

}

//return all 16 registers with value 0
func InitRegisters() []register {
	var regList []register
	for i := 0; i <= 15; i++ {
		regList = append(regList, register{number: i, value: 0})
	}
	return regList
}

func InitInAndOut() []inAndOut {
	var ioList []inAndOut
	for i := 0; i <= 1; i++ {
		ioList = append(ioList, inAndOut{number: i, value: 0})
	}
	return ioList

}

func load(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("LOAD Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing LOAD")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}

	a.setValue(regList[number].value)

	if debug {
		fmt.Println("LOAD Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func infunc(parameter string) {
	//check if parameter is 1 or 0 else the instruction is invalid
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("IN Could not convert", parameter)
		os.Exit(2)
	}
	if number != 1 && number != 0 {
		fmt.Println("IN Not a valid instruction")
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing IN")
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}

	//prompt user for input
	reader := bufio.NewReader(os.Stdin)
	var input string
	fmt.Println("Provide input:")
	input, _ = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)

	floatInput, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("IN Could not convert", input, err)
		os.Exit(2)
	}
	//put the input into the given io
	ioList[number].setValue(floatInput)
	a.setValue(ioList[number].value)
	if debug {
		fmt.Println("IN Done")
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}
	
}

func start() {

}

func stop() {
	if debug {
		fmt.Println("Executing STOP")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	fmt.Println("Program Terminated")
	os.Exit(0)
}


func out(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("OUT Could not convert", parameter)
		os.Exit(2)
	}
	if number != 1 && number != 0 {
		fmt.Println("OUT Not a valid instruction")
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing OUT")
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	ioList[number].setValue(a.value)
	fmt.Println("I/O", number, "has the value:", ioList[number].value)
	if debug {
		fmt.Println("OUT Done")
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}
}

func loadnum(parameter string) {
	number, err := strconv.ParseFloat(parameter, 64)
	if err != nil {
		fmt.Println("LOADNUM Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing LOADNUM")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	a.setValue(number)

	if debug {
		fmt.Println("LOADNUM Done")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func store(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("STORE Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing STORE")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}
	regList[number].setValue(a.value)
	if debug {
		fmt.Println("STORE Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}
}

func jumpneg(parameter string) (int, error) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("JUMPNEG Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing JUMPNEG")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	if math.Signbit(a.value) {
		if debug {
			fmt.Println("Jump gets executed")
		 }
		return number, nil
	} else {
		if debug {
			fmt.Println("Jump does not get executed")
		 }
		return 0, fmt.Errorf("number is not positive")
	}
}

func jumppos(parameter string) (int, error) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("JUMPPOS Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing JUMPPOS")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	if a.value == 0 {
		if debug {
			fmt.Println("Jump does not get executed")
		}
		return 0, fmt.Errorf("JUMPPOS number is zero")
	}

	if !math.Signbit(a.value) {
		if debug {
			fmt.Println("Jump gets executed")
		}
		return number, nil
	} else {
		if debug {
			fmt.Println("Jump does not get executed")
		}
		return 0, fmt.Errorf("JUMPPOS number is not positive")
	}

}

func jumpnull(parameter string) (int, error) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("JUMPNULL Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing JUMPNULL")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}

	if a.value != 0 {
		if debug {
			fmt.Println("Jump does not get executed")
		}
		return 0, fmt.Errorf("accumulator is not 0")
	} else {
		if debug {
			fmt.Println("Jump gets executed")
		}
		return number, nil
	}

}

func jump(parameter string) int {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("JUMP Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing JUMP")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	return number
}

func add(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("ADD Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing ADD")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := regList[number].value + a.value
	a.setValue(newValue)
	if debug {
		fmt.Println("ADD Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func addnum(parameter string) {
	number, err := strconv.ParseFloat(parameter, 64)
	if err != nil {
		fmt.Println("ADDNUM could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing ADDNUM")
		fmt.Println("Inhalt von Parameter is", number)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := number + a.value
	a.setValue(newValue)
	if debug {
		fmt.Println("ADDNUM DONE")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func sub(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("SUB Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing SUB")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value - regList[number].value
	a.setValue(newValue)
	if debug {
		fmt.Println("SUB Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func mul(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("MUL Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing MUL")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value * regList[number].value
	a.setValue(newValue)
	if debug {
		fmt.Println("MUL Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func div(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("DIV Could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing DIV")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value / regList[number].value
	a.setValue(newValue)
	if debug {
		fmt.Println("DIV Done")
		fmt.Println("Inhalt von Register", number, "is", regList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func subnum(parameter string) {
	number, err := strconv.ParseFloat(parameter, 64)
	if err != nil {
		fmt.Println("SUBNUM could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing SUBNUM")
		fmt.Println("Inhalt von Parameter is", number)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value - number
	a.setValue(newValue)
	if debug {
		fmt.Println("SUBNUM Done")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}

}

func mulnum(parameter string) {
	number, err := strconv.ParseFloat(parameter, 64)
	if err != nil {
		fmt.Println("MULNUM could not convert", parameter)
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing MULNUM")
		fmt.Println("Inhalt von Parameter is", number)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value * number
	a.setValue(newValue)
	if debug {
		fmt.Println("MULNUM Done")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}

func divnum(parameter string) {
	number, err := strconv.ParseFloat(parameter, 64)
	if err != nil {
		fmt.Println("DIVNUM could not convert", parameter)
		 os.Exit(2)
	}
	if debug {
		fmt.Println("Executing DIVNUM")
		fmt.Println("Inhalt von Parameter is", number)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	newValue := a.value / number
	a.setValue(newValue)
	if debug {
		fmt.Println("DIVNUM Done")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
}
