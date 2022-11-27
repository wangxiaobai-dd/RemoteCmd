package main

import (
	"RemoteCmd/Common"
	"bufio"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var regParam *regexp.Regexp

func init() {
	regParam = regexp.MustCompile("(BYTE|DWORD|WORD|char)\\s*(\\w+)\\[?([a-zA-Z0-9]*)\\]?")
}

type Server struct {
	Common.Server
}

func (s *Server) searchMessage(message string, messageFiles []string) (bool, url.Values) {

	str := "(MSG_BEGIN_?\\d*)\\([a-zA-Z0-9_]+\\s*,\\s*(\\d+)\\s*,\\s*" + message
	regBegin := regexp.MustCompile(str)

	response := url.Values{}
	findMsg := false

	for _, v := range messageFiles {
		f, err := os.Open(s.Path + "\\" + v)

		if err != nil {
			log.Println(err)
			continue
		}

		r := bufio.NewReader(f)
		for {
			line, err := r.ReadString('\n')

			ret := regBegin.FindStringSubmatch(line)
			if len(ret) != 0 {
				findMsg = true
				log.Println("find match:", ret)

				if len(ret) == 3 {
					log.Println("response set")
					response.Set("cmd", ret[1])
					response.Set("paraNum", ret[2])
				}
			}
			if findMsg {
				ret = regParam.FindStringSubmatch(line)
				if len(ret) == 4 {
					log.Println("find param match:", ret[3])
					response.Add(ret[2], ret[1]) // param type
					response.Add(ret[2], ret[3]) // array len
				}
				if strings.Contains(line, "MSG_END") {
					break
				}
			}
			if err != nil {
				break
			}
		}
		f.Close()

		if findMsg {
			break
		}
	}
	return findMsg, response
}
