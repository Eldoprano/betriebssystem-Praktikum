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
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments specified")
		return
	}
	d := flag.Bool("debug", false, "enables debug output")
	input := flag.String("input", "", "The HAL-OS Config File")
	flag.Parse()
	programs, cOne, cTwo, err := readConfFile(input)
	fmt.Println(programs)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var wg sync.WaitGroup
	buffer := make(chan string, 4)
	var blockInput int
	var blockOutput int
	fmt.Println(len(programs))
	for i := 0; i <= len(programs); i++ {
		if i == 0 {
			blockInput = 0
		} else {
			blockInput = cTwo[i]
		}
		blockOutput = cOne[i]
		wg.Add(1)
		go func() {
			m := readProgramFile(programs[i+1])
			fmt.Println(m)
			HAL.HalStart(m, *d, buffer, &wg, blockInput, blockOutput)
		}()
		wg.Wait()
	}
	fmt.Println("this is the end of the program!")
}

func readProgramFile(input string) map[int]string {
	var m map[int]string
	m = make(map[int]string)
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		lineNumber, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("here we are")
			log.Fatal("Line did not start with a number", err)
		}
		m[lineNumber] = strings.Join(s[1:], " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m

}


func readConfFile(input *string) (map[int]string, map[int]int, map[int]int, error) {
	fmt.Println("reading conf file")
	var e error
	var p map[int]string
	var cOne map[int]int
	var cTwo map[int]int
	p = make(map[int]string)
	cOne = make(map[int]int)
	cTwo = make(map[int]int)

	file, err := os.Open(*input)
	if err != nil {
		e = err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0
	connectionPart := false
	for scanner.Scan() {
		if counter == 0 && scanner.Text() != "HAL-Prozessoren:" {
				e = fmt.Errorf("Invalid Config File Format")
		}
		counter++
		if scanner.Text() == "HAL-Verbindungen:" {
			connectionPart = true

		}
		//lots of string formating to get the numbers from this horrible format
		if connectionPart == true && scanner.Text() != "HAL-Verbindungen:" {
			block1 := strings.Split(scanner.Text(), " > ")
			one, two := block1[0], block1[1]
			//mapone
			z := strings.Split(one, ":")
			mapValueOne, mapValueTwo := z[0], z[1]
			mapIntOne, err := strconv.Atoi(mapValueOne)
			if err != nil {
				e = fmt.Errorf("Invalid Format in Connections")
			}
			mapIntTwo, err := strconv.Atoi(mapValueTwo)
			if err != nil {
				e = fmt.Errorf("Invalid Format in Connections")
			}
			cOne[mapIntOne] = mapIntTwo

			//maptwo
			z = strings.Split(two, ":")
			mapValueOne, mapValueTwo = z[0], z[1]
			mapIntOne, err = strconv.Atoi(mapValueOne)
			if err != nil {
				e = fmt.Errorf("Invalid Format in Connections")
			}
			mapIntTwo, err = strconv.Atoi(mapValueTwo)
			if err != nil {
				e = fmt.Errorf("Invalid Format in Connections")
			}
			cTwo[mapIntOne] = mapIntTwo

		} else if scanner.Text() != "HAL-Prozessoren:" && scanner.Text() != "HAL-Verbindungen:" {
			s := strings.Fields(scanner.Text())
			lineNumber, err := strconv.Atoi(s[0])
			if err != nil {
				e = fmt.Errorf("Line did not start with a number, %s", s)
			}
			p[lineNumber] = strings.Join(s[1:], " ")
		} else {
			continue
		}
	}
	return p, cOne, cTwo, e
}