package main

import (
	"RemoteCmd/Common"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Server struct {
	Common.Server
}

var serverMap = make(map[string]Server)
var lock sync.Mutex
var proxy *httputil.ReverseProxy

func init() {
	remote, err := url.Parse(Common.ProxyUrl)
	if err != nil {
		log.Println(err)
		return
	}
	proxy = httputil.NewSingleHostReverseProxy(remote)
}

func main() {

	router := gin.Default()
	router.POST("/server/sync", serverSync)
	router.GET("/message/search/:serverName/:message", messageSearch)
	go checkServer()

	router.Run(Common.WorkerPort)
}

func serverSync(c *gin.Context) {
	buffer, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))

	server := Server{}
	err := c.ShouldBind(&server)
	if err != nil {
		log.Println("bind body err:", err)
		return
	}
	server.CheckTime = time.Now().Unix() + 2*60
	lock.Lock()
	serverMap[server.Name] = server
	lock.Unlock()
	log.Println("postServer:", server.Info())

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))

	forwardProxy(c)
}

func messageSearch(c *gin.Context) {
	serverName := c.Param("serverName")
	message := c.Param("message")
	log.Println("messageSearch,", serverName, "message:", message)

	// do search

	c.String(http.StatusOK, "OK")
}

func forwardProxy(c *gin.Context) {
	if proxy != nil {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func searchHeaderFile(f string) {
	_, err := os.ReadFile(f)
	if err != nil {
		return
	}
}

func checkServer() {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for {
			<-ticker.C
			lock.Lock()
			now := time.Now().Unix()
			for name, server := range serverMap {
				if server.CheckTime < now {
					log.Println("checkServer,delete", server.Info())
					deleteServerProxy(name)
					delete(serverMap, name)
				}
			}
			lock.Unlock()
		}
	}()
}

func deleteServerProxy(name string) {
	Common.SendRequest("DELETE", Common.ProxyUrl+"/server/delete/"+name, nil)
}
