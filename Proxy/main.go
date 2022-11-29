package main

import (
	"RemoteCmd/Common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var serverMap = make(map[string]Server)
var configMap = make(map[string]interface{})

type Server struct {
	Common.Server
}

func init() {
	data, err := os.ReadFile("Proxy/cmd.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, &configMap)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("configMap:", configMap)
	log.Println("messageFiles:", getMessageFiles())
}

func main() {
	router := gin.Default()
	router.POST("/server/sync", serverSync)
	router.DELETE("/server/delete/:serverName", serverDelete)
	router.GET("/message/search/:serverName/:message", messageSearch)
	router.POST("/message/send", messageSend)
	router.Run(Common.ProxyPort)
}

func serverSync(c *gin.Context) {
	server := Server{}
	err := c.ShouldBind(&server)
	if err != nil {
		log.Println("postServer,err:", err)
		return
	}
	if len(server.Name) == 0 {
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

	serverName := c.Param("serverName")
	server, ok := serverMap[serverName]
	if !ok {
		log.Println("messageSearch, no server:", serverName)
		return
	}

	request := url.Values{}
	request.Set("serverName", serverName)
	request.Set("message", c.Param("message"))
	for _, v := range getMessageFiles() {
		request.Add("messageFiles", v)
	}

	Common.SendRequestGin(c, "POST", server.UrlString()+"/message/search", strings.NewReader(request.Encode()))
}

func messageSend(c *gin.Context) {

	server, ok := serverMap[c.PostForm("serverName")]
	if !ok {
		c.JSON(http.StatusOK, gin.H{"status": "NoServer"})
		return
	}
	paramCount := c.PostForm("paramCount")
	log.Println("messageSend:", server.Info(), paramCount)

	// todo send http request or redis

}

func forwardWorker(server *Server, c *gin.Context) {
	remote, err := url.Parse(server.UrlString())
	if err != nil {
		return
	}
	worker := httputil.NewSingleHostReverseProxy(remote)
	worker.ServeHTTP(c.Writer, c.Request)
}

func getMessageFiles() []string {
	messageFiles, _ := configMap["MESSAGE_FILE"]
	var ret []string
	for _, v := range messageFiles.([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}
