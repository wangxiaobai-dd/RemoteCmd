package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var headers = [...]string{"cmd.h", "cmd2.h"}

func main() {

	router := gin.Default()
	router.POST("/serverInfo", postServerInfo)

	fmt.Println("worker")
	for _, header := range headers {
		dealHeaderFile(header)
	}
}

func postServerInfo(c *gin.Context) {

}

func dealHeaderFile(f string) {
	_, err := os.ReadFile(f)
	if err != nil {
		return
	}
}
