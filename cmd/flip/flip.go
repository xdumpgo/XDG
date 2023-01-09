package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("proxies-living.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	out, err := os.Create("proxies.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p := strings.Split(scanner.Text(), ":")
		out.WriteString(fmt.Sprintf("%s:%s:%s:%s\r\n", p[2], p[3], p[0], p[1]))
	}
}
