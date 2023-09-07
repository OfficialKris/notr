package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"
)

type InternalData struct {
	Week     int
	Classes  []Class
	LastDate time.Time
}

type Class struct {
	Code        string
	Description string
}

func main() {
	clearOutput()
	internalData := readInternalData()

	for {

		internalData.LastDate = time.Now()
		year, month, day := internalData.LastDate.Date()
		writeInternalData(internalData)

		fmt.Printf("Current Date: %v %v, %v\n", month, day, year)
		fmt.Printf("Current Week: %v\n", internalData.Week)

		fmt.Println()

		fmt.Println("a: Start a new NOTR")
		fmt.Println("b: Create a new class")
		fmt.Println("c: Change week number")
		fmt.Println("d: Change class information")
		fmt.Println("e: Set class directory")

		fmt.Println()

		fmt.Print("Enter menu option (a, b, or c, or q to quit): ")
		var first string

		fmt.Scan(&first)

		if first[0] == 'q' {
			return
		}

		switch first {
		case "a":
			createNOTR(internalData)
		case "b":
			//createClass
		case "c":
			//changeWeek
		case "d":
			//changeClassInformation
		case "e":
			//setClassDirectory
		}

	}
}

func createNOTR(internalData InternalData) {
	//fmt.Println(internalData.classes)
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
	fmt.Printf("Class %v\n", class)
	year, month, day := internalData.LastDate.Date()
	fmt.Printf("Current Date: %v-%v-%v\n", year, month, day)
	fmt.Printf("Current Week %v\n", internalData.Week)
	fmt.Printf("Started TextEdit with file name: %v\n", "NOTES")
	createFile("NOTES")
	stub := generateStub(class.Code, internalData.Week, internalData.LastDate)
	populateFile("NOTES", stub)
	openFile("NOTES")
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
	f, err := os.OpenFile("./"+name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
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
	path := usr.HomeDir + "/.config/notr/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	err1 := os.WriteFile(path+"config.json", []byte(jsonData), 0644)
	check(err1)
}

func readInternalData() (returnData InternalData) {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.config/notr/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}
	data, err := os.ReadFile(path + "config.json")
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
