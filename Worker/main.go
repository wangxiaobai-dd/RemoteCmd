package main

import (
	"os"
)

var headers = [...]string{"cmd.h", "cmd2.h"}

func main() {

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
