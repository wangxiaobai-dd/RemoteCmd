package main

import (
	"RemoteCmd/Common"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var regParam *regexp.Regexp

func init() {
	regParam = regexp.MustCompile("(BYTE|WORD|DWORD|QWORD|char)\\s*(\\w+)\\[?([a-zA-Z0-9]*)\\]?")
}

type Server struct {
	Common.Server
}

func (s *Server) searchMessage(message string, messageFiles []string) map[string]interface{} {

	str := "(MSG_BEGIN_?\\d*)\\([a-zA-Z0-9_]+\\s*,\\s*(\\d+)\\s*,\\s*" + message
	regBegin := regexp.MustCompile(str)

	response := make(map[string]interface{})
	paramCount := 0

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
				paramCount = 1
				log.Println("find match:", ret)

				if len(ret) == 3 {
					log.Println("response set")
					response["cmdFlag"] = ret[1]
					response["paraNumber"] = ret[2]
				}
			}
			if paramCount > 0 {
				ret = regParam.FindStringSubmatch(line)
				if len(ret) == 4 {
					// ret[1]:paramName ret[2]:type ret[3]:array len
					log.Println("find param match:", ret[1], ret[2], ret[3])

					paramKey := fmt.Sprint("param", paramCount)

					paramArr := []string{ret[1], ret[2], ret[3]}
					response[paramKey] = paramArr
					paramCount++
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

		if paramCount > 0 {
			response["paramCount"] = paramCount - 1
			break
		}
	}

	return response
}
