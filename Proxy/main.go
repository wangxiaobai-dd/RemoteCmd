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
	router.StaticFS("/plugin", http.Dir("Web/plugin"))
	router.StaticFS("/static", http.Dir("Web/static"))
	router.LoadHTMLFiles(FrontDir + "main.html")

	router.GET("/", pageShow)
	router.GET("/server/showList", serverShowList)
	router.GET("/user/showList/", userShowList)

	router.POST("/server/sync", serverSync)
	router.DELETE("/server/delete/:serverName", serverDelete)
	router.GET("/message/search/", messageSearch)
	router.POST("/message/send", messageSend)

	router.Run(Common.ProxyPort)
}

func pageShow(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{})
}

func serverShowList(c *gin.Context) {
	var response []string
	for serverName, _ := range serverMap {
		response = append(response, serverName)
	}

	c.JSON(http.StatusOK, response)
}

func userShowList(c *gin.Context) {
	serverName, _ := c.GetQuery("serverName")
	log.Println("userShowList:", serverName)
	server, ok := serverMap[serverName]
	if !ok {
		return
	}
	var response []string
	for _, user := range server.Users {
		response = append(response, user)
	}
	c.JSON(http.StatusOK, response)
}

func serverSync(c *gin.Context) {
	server := Server{}
	//err := c.ShouldBindJSON(&server) //json
	err := c.ShouldBind(&server) //urlencoded
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

	serverName, _ := c.GetQuery("serverName")
	message, _ := c.GetQuery("message")
	server, ok := serverMap[serverName]
	if !ok {
		log.Println("messageSearch, no server:", serverName)
		return
	}

	request := url.Values{}
	request.Set("serverName", serverName)
	request.Set("message", message)
	for _, v := range getMessageFiles() {
		request.Add("messageFiles", v)
	}
	log.Println("messageSearch:", serverName, message)

	Common.SendRequestGin(c, "POST", server.UrlString()+"/message/search", strings.NewReader(request.Encode()), requestRetCBFunc)
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
	log.Println("messageSend:", body)

	rdb := server.getRedisDb()
	log.Println(rdb.RPush(server.Path+RedisKey, data).Result())
}

func getMessageFiles() []string {
	messageFiles, _ := configMap["MESSAGE_FILE"]
	var ret []string
	for _, v := range messageFiles.([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}

func requestRetCBFunc(reqUrl string, bodyRet *[]byte) {

	if strings.Contains(reqUrl, "message/search") {
		retMap := make(map[string]interface{})
		json.Unmarshal(*bodyRet, &retMap)

		cmdFlag, ok := retMap["cmdFlag"]
		if !ok {
			return
		}
		cmdNumber, ok := configMap[cmdFlag.(string)]
		if ok {
			retMap["cmdNumber"] = cmdNumber.(string)
			*bodyRet, _ = json.Marshal(retMap)
		}
	}
}

func forwardWorker(server *Server, c *gin.Context) {
	remote, err := url.Parse(server.UrlString())
	if err != nil {
		return
	}
	worker := httputil.NewSingleHostReverseProxy(remote)
	worker.ServeHTTP(c.Writer, c.Request)
}
