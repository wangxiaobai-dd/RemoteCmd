package main

import (
	"RemoteCmd/Common"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http/httputil"
)

type Server struct {
	Common.Server
}

var serverMap = make(map[string]Server)

func main() {
	router := gin.Default()
	router.POST("/server/sync", serverSync)
	router.DELETE("/server/delete/:serverName", serverDelete)
	router.GET("/message/search/:serverName/:message", messageSearch)
	router.Run(Common.ProxyPort)
}

func serverSync(c *gin.Context) {
	server := Server{}
	err := c.ShouldBind(&server)
	if err != nil {
		log.Println("postServer,err:", err)
		return
	}
	serverMap[server.Name] = server
	log.Println("serverSync:", server.Info())
}

func serverDelete(c *gin.Context) {
	serverName := c.Param("serverName")
	delete(serverMap, serverName)
	log.Println("deleteServer:", serverName, "len:", len(serverMap))
}

func messageSearch(c *gin.Context) {

	buffer, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))

	serverName := c.Param("serverName")
	server, ok := serverMap[serverName]
	if !ok {
		log.Println("messageSearch, no server:", serverName)
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	forwardWorker(&server, c)
}

func forwardWorker(server *Server, c *gin.Context) {
	worker := httputil.NewSingleHostReverseProxy(server.Url())
	worker.ServeHTTP(c.Writer, c.Request)
}
