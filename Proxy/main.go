package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const jsonFile = "Proxy/server.json"
const remotePort = "7001"

var userMap = make(map[string]interface{})

func init() {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, &userMap)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(userMap)
}

func main() {
	router := gin.Default()
	router.GET("/:name", proxy)
	router.Run(":8080")
}

func proxy(c *gin.Context) {
	c.String(http.StatusOK, "HelloWorld")
	name := c.Param("name")
	ip, ok := userMap[name]
	if !ok {
		log.Println("no user:", name)
		return
	}

	log.Println("proxy,found user:", name, "ip:", ip)

	remote, err := url.Parse("http://" + ip.(string) + ":" + remotePort)
	if err != nil {
		log.Println(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}
