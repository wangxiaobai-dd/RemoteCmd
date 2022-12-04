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

const (
	RedisKey = ":mock.message"
	CmdJson  = "Proxy/cmd.json"
	FrontDir = "Web/"
)

func init() {
	data, err := os.ReadFile(CmdJson)
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

	router.Delims("[[", "]]")
	router.StaticFS("/css", http.Dir("Web/css"))
	router.StaticFS("/js", http.Dir("Web/js"))
	router.LoadHTMLFiles(FrontDir + "main.html")

	router.GET("/", pageShow)
	router.GET("/server", getServerList)

	router.POST("/server/sync", serverSync)
	router.DELETE("/server/delete/:serverName", serverDelete)
	router.GET("/message/search/:serverName/:message", messageSearch)
	router.POST("/message/send", messageSend)

	router.Run(Common.ProxyPort)
}

func pageShow(c *gin.Context) {

	c.HTML(http.StatusOK, "main.html", gin.H{})
}

func getServerList(c *gin.Context) {
	// todo server list
}

func serverSync(c *gin.Context) {
	server := Server{}
	err := c.ShouldBind(&server)
	if err != nil {
		log.Println("postServer,err:", err)
		return
	}
	if len(server.ServerName) == 0 {
		return
	}
	serverMap[server.ServerName] = server
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

	body := make(map[string]interface{})
	data, _ := c.GetRawData()
	err := json.Unmarshal(data, &body)
	if err != nil {
		log.Println("messageSend,err:", err)
		return
	}
	server, ok := serverMap[body[("serverName")].(string)]
	if !ok {
		c.JSON(http.StatusOK, gin.H{"status": "NoServer"})
		return
	}
	log.Println("messageSend:", server.Info())
	rdb := server.getRedisDb()
	c.JSON(http.StatusOK, body)
	log.Println(rdb.RPush(server.Path+RedisKey, body).Result())
}

func getMessageFiles() []string {
	messageFiles, _ := configMap["MESSAGE_FILE"]
	var ret []string
	for _, v := range messageFiles.([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}

func forwardWorker(server *Server, c *gin.Context) {
	remote, err := url.Parse(server.UrlString())
	if err != nil {
		return
	}
	worker := httputil.NewSingleHostReverseProxy(remote)
	worker.ServeHTTP(c.Writer, c.Request)
}
