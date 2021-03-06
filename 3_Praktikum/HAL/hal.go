package HAL

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	pc              int = 0
	regList         []register
	ioList          []inAndOut
	a               accu
	debug           = false
	instructionList map[int]string
	WaitGroup       *sync.WaitGroup
)

type Connection struct {
	Port     int
	Channel  chan float64
	ConnType string
}

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

// connection : verbindung []map[int]map[string]int
func closeChannels(connectionList []Connection, prozessorNum int) {
	for _, conn := range connectionList {
		if conn.ConnType == "from" {
			//fmt.Println(prozessorNum, "is closing channel", conn.Channel, "from list", connectionList)
			close(conn.Channel)
		}
	}
	WaitGroup.Done()
}

func HalStart(instr map[int]string, d bool, num int, conn []Connection, wg *sync.WaitGroup) {
	if d {
		debug = true
		fmt.Println("Running in Debug Mode")
	}
	//initalize all stuff
	instructionList = instr
	regList = InitRegisters()
	ioList = InitInAndOut(20) // Yeah... 20 👀
	a = accu{value: 0}
	WaitGroup = wg
	stopCalled := false
	defer closeChannels(conn, num)

	//parse instruction and parameters
	for i := 1; i <= len(instructionList); i++ {
		if debug {
			time.Sleep(0 * time.Second)
			fmt.Println("-----------------------------------")
		}
		value := instructionList[i]
		s := strings.Fields(value)
		parameter := strings.Join(s[1:], " ")
		//fmt.Println(parameter)
		switch instruction := s[0]; instruction {
		case "START":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
				fmt.Println("Starting")
			}
			//start(parameter)
		case "STOP":
			stopCalled = true
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			stop(num)
		case "OUT":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			out(parameter, num, conn)
		case "IN":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			infunc(parameter, num, conn)
		case "LOAD":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			load(parameter)
		case "LOADNUM":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			loadnum(parameter)
		case "STORE":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			store(parameter)
		case "JUMPNEG":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			address, err := jumpneg(parameter)
			if err != nil {
				if debug {
					//fmt.Println("not jumping")
				}
			} else {
				i = address - 1
			}
		case "JUMPPOS":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			address, err := jumppos(parameter)
			if err != nil {
				//fmt.Println("not jumping")
			} else {
				i = address - 1
			}
		case "JUMPNULL":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			address, err := jumpnull(parameter)
			if err != nil {
				//fmt.Println("not jumping")
			} else {
				i = address - 1
			}
		case "JUMP":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			address := jump(parameter)
			i = address - 1
		case "ADD":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			add(parameter)
		case "ADDNUM":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			addnum(parameter)
		case "SUB":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			sub(parameter)
		case "MUL":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			mul(parameter)
		case "DIV":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			div(parameter)
		case "SUBNUM":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			subnum(parameter)
		case "MULNUM":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			mulnum(parameter)
		case "DIVNUM":
			if debug {
				fmt.Println("-----------Prozessor ", num, "-----------")
			}
			divnum(parameter)
		default:
			fmt.Println("Invalid Instruction", instruction)
			os.Exit(2)
		}
		if stopCalled {
			break
		}
		//fmt.Println("Num:", num, "Connections:", conn)
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

func InitInAndOut(numberOfIO int) []inAndOut {
	var ioList []inAndOut
	for i := 0; i < numberOfIO; i++ {
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

func infunc(parameter string, prozessorNum int, connectionList []Connection) {
	//check if parameter is 1 or 0 else the instruction is invalid
	number, err := strconv.Atoi(parameter)
	var floatInput float64
	if err != nil {
		fmt.Println("IN Could not convert", parameter)
		os.Exit(2)
	}
	if number >= len(ioList) || number < 0 {
		fmt.Println("IN Not a valid instruction")
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing IN Prozessor", prozessorNum)
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}

	receivedPerChannel := false
	for _, conn := range connectionList {
		if conn.ConnType == "to" && number == conn.Port {
			//fmt.Println("INCHANNEL Prozessor", prozessorNum, "is listening channel:", conn.Channel)
			floatInput = <-conn.Channel
			receivedPerChannel = true
			//fmt.Println("INCHANNEL executed Prozessor", prozessorNum, "Received numer", floatInput)
		}
	}

	if !receivedPerChannel {
		//prompt user for input
		reader := bufio.NewReader(os.Stdin)
		var input string
		fmt.Println("Provide input:")
		input, _ = reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		floatInput, err = strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("IN Could not convert", input, err)
			os.Exit(2)
		}
	}
	//put the input into the given io
	ioList[number].setValue(floatInput)
	a.setValue(ioList[number].value)
	if debug {
		fmt.Println("IN Done Prozessor", prozessorNum)
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhatl von Akkumulator ist: ", a.value)
	}

}

func start() {

}

func stop(prozessorNum int) {
	if debug {
		fmt.Println("Executing STOP")
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}
	fmt.Println("Prozessor", prozessorNum, "Terminated")
}

func out(parameter string, prozessorNum int, connectionList []Connection) {
	number, err := strconv.Atoi(parameter)
	if err != nil {
		fmt.Println("OUT Could not convert", parameter)
		os.Exit(2)
	}
	if number >= len(ioList) || number < 0 {
		fmt.Println("OUT Not a valid instruction")
		os.Exit(2)
	}
	if debug {
		fmt.Println("Executing OUT Prozessor", prozessorNum)
		fmt.Println("Inhalt von I/O", number, "is", ioList[number].value)
		fmt.Println("Inhalt von Akkumulator ist: ", a.value)
	}

	//fmt.Println("Prozessor:", prozessorNum, " is Searching in: ", connectionList, " The number: ", number)
	for _, conn := range connectionList {
		if conn.ConnType == "from" && number == conn.Port {
			conn.Channel <- a.value
			if debug {
				fmt.Println("OUTCHANNEL Done Prozessor", prozessorNum, "channel:", conn.Channel)
				fmt.Println("Inhalt von I/O", number, "is", a.value)
				fmt.Println("Inhatl von Akkumulator ist: ", a.value)
			}
			return
		}
	}

	ioList[number].setValue(a.value)
	fmt.Println("On Prozessor", prozessorNum, " the I/O", number, "has the value:", ioList[number].value)
	if debug {
		fmt.Println("OUT Don Prozessor", prozessorNum)
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

	// Signbit reports whether x is negative or negative zero.
	if math.Signbit(a.value) {
		if debug {
			fmt.Println("Jump gets executed")
		}
		return number, nil
	}
	if debug {
		fmt.Println("Jump does not get executed")
	}
	return 0, fmt.Errorf("number is not positive") //Should we really use Errorf?
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
	}
	if debug {
		fmt.Println("Jump does not get executed")
	}
	return 0, fmt.Errorf("JUMPPOS number is not positive")
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
	}
	if debug {
		fmt.Println("Jump gets executed")
	}
	return number, nil

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
