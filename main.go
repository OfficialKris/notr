package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/officialkris/notr/config"
)

var DevMode = false

type InternalData struct {
	Version        string
	Week           int
	Classes        []Class
	LastDate       time.Time
	ClassDirectory string
}

type Class struct {
	Code        string
	Description string
}

func main() {
	clearOutput()
	internalData := readInternalData()

	checkData(&internalData)

	for {
		internalData.Version = config.PROGRAM_VERSION
		internalData.LastDate = time.Now()
		year, month, day := internalData.LastDate.Date()
		writeInternalData(internalData)

		if config.DEV_MODE {
			fmt.Println("DEV_MODE ENABLED: USE WITH CARE")
		}

		fmt.Printf("NOTR Version: %v\n\n", internalData.Version)
		fmt.Printf("Current Date: %v %v, %v\n", month, day, year)
		fmt.Printf("Current Week: %v\n", internalData.Week)
		fmt.Printf("Current Config Directory: ~%v\n", config.CONFIG_LOCATION+config.CONFIG_FILE_NAME)
		fmt.Printf("Current Directory: %v\n", internalData.ClassDirectory)

		fmt.Println()

		fmt.Println("a: Start a new NOTR")
		fmt.Println("b: Modify classes")
		fmt.Println("c: Change week number")
		fmt.Println("d: Move to new Semester")
		fmt.Println("e: Set class directory")

		fmt.Println()

		fmt.Print("Enter menu option (a to e, or q to quit): ")
		var first string

		fmt.Scan(&first)

		first = strings.ToLower(first)

		if first[0] == 'q' {
			return
		}

		switch first {
		case "a":
			createNOTR(internalData)
		case "b":
			//modifyClasses
			notImplemented()
		case "c":
			setWeek(&internalData)
		case "d":
			//newSemester
			notImplemented()
		case "e":
			setClassDirectory(&internalData)
		}

	}
}

func setWeek(internalData *InternalData) {
	for {
		fmt.Println()
		fmt.Print("Enter new week number (or q to quit): ")
		var first string
		fmt.Scan(&first)

		if first[0] == 'q' {
			return
		}

		num, err := strconv.Atoi(first)
		if err != nil {
			fmt.Println("Please enter a valid number.")
		} else {
			internalData.Week = num
			fmt.Println()
			break
		}
	}
}

func setClassDirectory(internalData *InternalData) {
	for {
		fmt.Println()
		fmt.Print("Enter fully resolved path name for new directory (or q for quit): ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}

		line = strings.TrimRight(line, "\n")

		if line[0] == 'q' {
			return
		}

		if _, err := os.Stat(line); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Folder does not exist.")
				fmt.Println(line)
			} else {
				panic(err)
			}
		} else {
			line, _ = filepath.Abs(line)
			fmt.Println("Folder DOES exist!")
			internalData.ClassDirectory = line
			break
		}
	}
}

func createNOTR(internalData InternalData) {
	fmt.Println()
	for i, s := range internalData.Classes {
		fmt.Printf("%v: %v - %v\n", i+1, s.Code, s.Description)
	}

	fmt.Println("q: Return to base menu")
	fmt.Println()

	fmt.Print("Enter class (from range above): ")
	var first string

	fmt.Scan(&first)

	if first[0] == 'q' {
		return
	}
	parsedInput, err := strconv.Atoi(first)
	if err != nil {
		panic(err)
	}
	class := internalData.Classes[parsedInput-1]
	year, month, day := internalData.LastDate.Date()
	name := generateNoteName()
	stub := generateStub(class.Code, internalData.Week, internalData.LastDate)
	filepath := internalData.ClassDirectory + "/" + name

	fmt.Println()
	fmt.Printf("Class: %v\n", class)
	fmt.Printf("Current Date: %v-%v-%v\n", year, month, day)
	fmt.Printf("Current Week %v\n", internalData.Week)
	fmt.Printf("Started TextEdit with file name: %v\n", name)

	createFile(filepath)
	populateFile(filepath, stub)
	openFile(filepath)

	fmt.Println()
}

func generateNoteName() string {
	return "NOTES" + ".txt" // TODO: Add number
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func generateStub(class string, week int, time time.Time) (buf string) {
	day := fmt.Sprintf("%s", time.Weekday())
	day = day[0:3]
	buf = fmt.Sprintf("%s\nWEEK %d\n%s %d/%d/%d\n\n", class, week, day, time.Month(), time.Day(), time.Year())
	return
}

func populateFile(name string, stub string) {
	fmt.Println(name)
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	check(err)
	defer f.Close()
	n2, err := f.Write([]byte(stub))
	check(err)
	fmt.Printf("Wrote bytes: %v", n2)
}

func createFile(name string) {
	app := "touch"

	cmd := exec.Command(app, name)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

func openFile(name string) {
	app := "open"
	arg1 := "-a"
	arg2 := "TextEdit"
	arg3 := name

	cmd := exec.Command(app, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

func writeInternalData(data InternalData) {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf("Could not marshal data: %s\n", err)
		return
	}
	usr, _ := user.Current()
	path := usr.HomeDir + config.CONFIG_LOCATION
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	err1 := os.WriteFile(path+config.CONFIG_FILE_NAME, []byte(jsonData), 0644)
	check(err1)
}

func readInternalData() (returnData InternalData) {
	usr, _ := user.Current()
	path := usr.HomeDir + config.CONFIG_LOCATION
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	data, err := os.ReadFile(path + config.CONFIG_FILE_NAME)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	json.Unmarshal(data, &returnData)
	return
}

func clearOutput() {
	cmd := exec.Command("clear")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

func notImplemented() {
	println()
	println("Function is not implemented yet")
	println()
}

func checkData(internalData *InternalData) {
	if internalData.ClassDirectory == "" {
		println("Unfortunately, the config file under ~/.config/notr does not have a value for the class directory.")
		setClassDirectory(internalData)
	}
}
