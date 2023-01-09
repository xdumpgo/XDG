package main

import (
	"bufio"
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/sqweek/dialog"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Simple password filter by Zertex#0001")
	fmt.Println("Select your combo file")

	filename, err := dialog.File().Filter("Combo list", "txt").Title("Select your combo file").Load()
	if err != nil {
		panic(err.Error())
	}

	f, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}

	mapCounter := make(map[string]int)
	var comboArray []string

	fmt.Print("Loading combos...  ")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		comboArray = append(comboArray, scanner.Text())
	}
	fmt.Println("Done!  Loaded",len(comboArray))

	fmt.Println("How many dupes would you like to set as max? [100]")
	var dupes = 100
	for {
		dupes = GetInputInt()
		if dupes <= 0 {
			fmt.Println("Please choose a value higher than 0")
		}
		break
	}

	bar := uiprogress.AddBar(len(comboArray))

	fmt.Println("Processing lines..")
	for _, comboline := range comboArray {
		s := strings.Split(comboline, ":")
		if len(s) == 2 {
			if _, ok := mapCounter[s[1]]; ok {
				mapCounter[s[1]]++
			} else {
				mapCounter[s[1]] = 1
			}
		}
		bar.Incr()
	}

	bar2 := uiprogress.AddBar(len(comboArray))
	output, err := os.Create("output.txt")
	if err != nil {
		panic(err.Error())
	}
	good := 0

	fmt.Println("Filtering and saving lines...")
	for _, comboline := range comboArray {
		s := strings.Split(comboline, ":")
		if val, ok := mapCounter[s[1]]; ok && val < dupes {
			good++
			output.WriteString(comboline + "\r\n")
		}
		bar2.Incr()
	}
	fmt.Println("Done! Saved",good,"lines to output.txt")
}
var inputReader = bufio.NewReader(os.Stdin)
func GetInputInt() int {
	fmt.Print("> ")
	s, err := inputReader.ReadString('\n')
	fmt.Println(s)
	if err != nil {
		return -1
	}
	i, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", ""))
	fmt.Println(i)
	if err != nil {
		return -1
	}
	return i
}