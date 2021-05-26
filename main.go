package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	MConsole = 0
	MFile    = 1
)

func main() {
	params := os.Args[1:]
	mode, input, output := getAppParams(params)
	switch mode {
	case MConsole:
		echo("calculate console")
		calcConsole()
		break
	case MFile:
		echo("calculate file")
		calcFile(input, output)
		break
	default:
		log.Fatalf("unknow app mode %v. stop app working. app arguments: %v\n", mode, params)
	}
	fmt.Println("all calcs have been completed")
}

func getAppParams(params []string) (mode int, input string, output string) {
	mode = MConsole
	if len(params) != 2 {
		return mode, input, output
	}

	input, output = params[0], params[1]
	if input != "-" && fileExists(input) {
		mode = MFile
	}
	return mode, input, output
}

func calcConsole() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error while reading console %v\n", err)
		}
		if strings.Contains(strings.ToUpper(line), "EXIT") {
			break
		}

		result := calcLine(line)
		printResultOnConsole(result)
	}
}

func calcFile(input, output string) {
	file, err := os.Open(input)
	if err != nil {
		log.Fatalf("error while reading file %v %v\n", input, err)
	}
	defer func() { _ = file.Close() }()

	if fileExists(output) {
		err := os.Remove(output)
		if err != nil {
			log.Fatalf("error while deleting file %v %v\n", output, err)
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		echo("->", line)
		result := calcLine(line)
		printResultOnFile(result, output)
	}
}

func calcLine(line string) interface{} {
	// TODO do calculation
	line = strings.Replace(line, "\n", "", -1)
	params := strings.Split(line, " ")
	const cmdSize = 3
	if len(params) < cmdSize || ((len(params)-cmdSize)%2) != 0 {
		return fmt.Errorf("error: incorrect format. (examples 'add 1 2' or 'mul 1 2 add 1' (without quotes))")
	}

	err, result := calc(params[0], params[1], params[2])
	if err != nil {
		return err
	}

	for i := 0; i < (len(params) - cmdSize); i = i + 2 {
		cmd := params[cmdSize:][i]
		op2 := params[cmdSize:][i+1]
		op1 := fmt.Sprintf("%f", result)

		err, result = calc(cmd, op1, op2)
		if err != nil {
			return err
		}
	}

	return result
}

func calc(command string, operand1, operand2 string) (error, float64) {
	var left, right float64
	var err error

	if left, err = strconv.ParseFloat(operand1, 64); err != nil {
		return fmt.Errorf("error while converting %v to float %v", operand1, err), 0
	}
	if right, err = strconv.ParseFloat(operand2, 64); err != nil {
		return fmt.Errorf("error while converting %v to float %v", operand2, err), 0
	}

	switch strings.ToUpper(command) {
	case "ADD":
		return nil, left + right
	case "MUL":
		return nil, left * right
	case "SUB":
		return nil, left - right
	case "DIV":
		return nil, left / right
	default:
		return fmt.Errorf("command %v is not supported", command), 0
	}
}

func echo(v ...interface{}) {
	fmt.Println(v...)
}

func printResultOnConsole(result interface{}) {
	echo("->", "answer is", result)
}

func printResultOnFile(result interface{}, output string) {
	file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		echo("error while opening or creating output file", output, err)
		printResultOnConsole(result)
		return
	}
	defer func() { _ = file.Close() }()

	if _, err := file.WriteString(fmt.Sprintf("%v\n", result)); err != nil {
		echo("error while writing result into file", output, err)
	}
	printResultOnConsole(result)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
