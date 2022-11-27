package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ProxyIp    = "127.0.0.1"
	ProxyPort  = ":7000"
	WorkerPort = ":7001"
)

var ProxyUrl = "http://" + ProxyIp + ProxyPort

type Server struct {
	Ip        string   `form:"ip" json:"ip"`
	Name      string   `form:"name" json:"name"`
	Path      string   `form:"path" json:"path"`
	Users     []string `form:"users" json:"users"`
	CheckTime int64
}

func (s *Server) Info() string {
	ret := "Name:" + s.Name + ",Ip:" + s.Ip + ",Path:" + s.Path + ",Users:"
	for _, user := range s.Users {
		ret = ret + user + "|"
	}
	return ret
}

func (s *Server) UrlString() string {
	return "http://" + s.Ip + WorkerPort
}

func SendRequest(method string, url string, body io.Reader) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SendRequestGin(c *gin.Context, method string, url string, body io.Reader) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyRet, _ := ioutil.ReadAll(resp.Body)

	c.String(http.StatusOK, string(bodyRet))
}
