package Common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (s *Server) Url() *url.URL {
	remote, _ := url.Parse("http://" + s.Ip + WorkerPort)
	return remote
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
