package main

import (
	"RemoteCmd/Common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct {
	Common.Server
}

var serverMap = make(map[string]Server)

func init() {

}

func main() {
	router := gin.Default()
	router.POST("/postServer", postServer)
	router.DELETE("/deleteServer", deleteServer)
	router.Run(Common.ProxyPort)
}

func postServer(c *gin.Context) {

	server := Server{}
	server.Ip = c.PostForm("ip")
	server.Name = c.PostForm("name")
	server.Path = c.PostForm("path")
	server.Users = c.PostFormArray("users")
	serverMap[server.Name] = server
	log.Println("postServer:", server.Info())
	c.String(http.StatusOK, server.Info())
}

func deleteServer(c *gin.Context) {
	delete(serverMap, c.PostForm("name"))
	log.Println("deleteServer:", c.PostForm("name"))
}
