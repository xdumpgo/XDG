package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Invalid arguments, please specify a file")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	outFile, err := os.Create(fmt.Sprintf("%s-fixed.txt", strings.Split(f.Name(), ":")[0]))
	if err != nil {
		log.Fatal(err.Error())
	}

	for scanner := bufio.NewScanner(f);scanner.Scan(); {
		var sep string
		if strings.Contains(scanner.Text(), ":") {
			sep = ":"
		} else {
			sep = ","
		}
		s := strings.Split(scanner.Text(), sep)
		fmt.Println(sep,s)
		outFile.WriteString(fmt.Sprintf("%s:%s\r\n", s[1], s[0]))
	}
	fmt.Println("Done!")
}
