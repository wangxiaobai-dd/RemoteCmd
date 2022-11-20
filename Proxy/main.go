package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type User struct {
	Name string
	Ip   string
}

const jsonFile = "Proxy/server.json"

var users []User

func init() {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(users)
}

func main() {

	r := gin.Default()

	r.Any("/*proxyPath", proxy)
}

func proxy(c *gin.Context) {

}
