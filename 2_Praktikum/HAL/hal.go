package HAL

import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"bufio"
	"os"
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
var pc int = 0
var regList []register
var ioList []inAndOut
var a accu

func HalStart(instr map[int]string) {
	fmt.Println("starting...")

	//initalize all stuff
	regList = InitRegisters()
	ioList = InitInAndOut()
	a = accu{value: 0}

	//parse instruction and parameters
	for _, value := range instr {
		s := strings.Fields(value)
		parameter := strings.Join(s[1:], " ")
		//fmt.Println(parameter)
		switch instruction := s[0]; instruction {
		case "Start":
			//start(parameter)
		case "Stop":
			//stop(parameter)
		case "OUT":
			//out(parameter)
		case "IN":
			infunc(parameter)
		case "LOAD":
			//load(parameter)
		case "LOADNUM":
			//loadnum(parameter)
		case "STORE":
			//store(parameter)
		case "JUMPNEG":
			//jumpneg(parameter)
		case "JUMPPOS":
			//jumppos(parameter)
		case "JUMPNULL":
			//jumpnull(parameter)
		case "JUMP":
			//jump(parameter)
		case "ADD":
			//add(parameter)
		case "ADDNUM":
			//addnum(parameter)
		case "SUB":
			//sub(parameter)
		case "MUL":
			//mul(parameter)
		case "DIV":
			//div(parameter)
		case "SUBNUM":
			//subnum(parameter)
		case "MULNUM":
			//mulnum(parameter)
		case "DIVNUM":
			//divnum(parameter)
		default:
			fmt.Println("Invalid Instruction")
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

func getInput() float64 {
	return 5
}


func load() {
	for _, reg2 := range regList {
		fmt.Println(reg2.number, math.Round(reg2.value))
	}
}

func infunc(parameter string) {
	//check if parameter is 1 or 0 else the instruction is invalid
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("Could not convert", parameter)
	}
	if number != 1 && number != 0 {
		fmt.Println("Not a valid instruction")
	}

	//prompt user for input
	reader := bufio.NewReader(os.Stdin)
	var input string
	fmt.Println("Provide input:")
	input, _ = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)

	floatInput, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Could not convert", input, err)
	}
	//put the input into the given io
	ioList[number].setValue(floatInput)
	a.setValue(ioList[number].value)
	fmt.Println(a.value)
	
}

func start() {

}

func stop() {

}

func out(parameter string) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("Could not convert", parameter)
	}
	if number != 1 && number != 0 {
		fmt.Println("Not a valid instruction")
	}


}

func loadnum() {

}

func store() {

}

func jumpneg() {

}

func jumppos() {

}

func jumpnull() {

}

func jump() {

}

func add() {

}

func addnum() {

}

func sub() {

}

func mul() {

}

func div() {

}

func subnum() {

}

func mulnum() {

}

func divnum() {

}