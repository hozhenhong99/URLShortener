package cache

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Cache map[string]string

func InitCache() {
	Cache = make(map[string]string)
	wd, _ := os.Getwd()
	file, err := os.Open(wd + "/cache/db")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			Cache[key] = value
		} else {
			fmt.Println("Invalid line:", line)
		}
	}
}

func WriteToDb(key, value string) {
	wd, _ := os.Getwd()
	filePath := wd + "/cache/db"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s %s\n", key, value))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
