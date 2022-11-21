package main

import (
	"fmt"
	"os"
)

var headers = [...]string{"cmd.h", "cmd2.h"}

func main() {
	fmt.Println("worker")
	for _, header := range headers {
		dealHeaderFile(header)
	}
}

func dealHeaderFile(f string) {
	_, err := os.ReadFile(f)
	if err != nil {
		return
	}
}
