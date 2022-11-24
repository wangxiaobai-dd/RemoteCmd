package main

import (
	"RemoteCmd/Common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Server struct {
	Common.Server
}

var headers = [...]string{"cmd.h", "cmd2.h"}
var serverMap = make(map[string]Server)
var lock sync.Mutex
var proxy *httputil.ReverseProxy

func init() {
	remote, err := url.Parse("http://" + Common.ProxyIp + Common.ProxyPort)
	if err != nil {
		log.Println(err)
		return
	}
	proxy = httputil.NewSingleHostReverseProxy(remote)
}

func main() {

	router := gin.Default()
	router.POST("/postServer", postServer)

	go checkServer()

	router.Run(Common.WorkerPort)
}

func postServer(c *gin.Context) {
	server := Server{}
	server.Ip = c.PostForm("ip")
	server.Name = c.PostForm("name")
	server.Path = c.PostForm("path")
	//server.Users = c.PostFormArray("users")
	server.CheckTime = time.Now().Unix() + 2*60
	lock.Lock()
	serverMap[server.Name] = server
	lock.Unlock()
	log.Println("postServer:", server.Info())

	forward(c)
}

func forward(c *gin.Context) {
	if proxy != nil {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func dealHeaderFile(f string) {
	_, err := os.ReadFile(f)
	if err != nil {
		return
	}
}

func checkServer() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	go func() {
		for {
			<-ticker.C
			lock.Lock()
			now := time.Now().Unix()
			for name, server := range serverMap {
				if server.CheckTime < now {
					log.Println("checkServer,delete", name)
					delete(serverMap, name)
				}
			}
			lock.Unlock()
		}
	}()
}
