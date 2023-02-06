package main

import (
	"RemoteCmd/Common"
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var regParam *regexp.Regexp

func init() {
	regParam = regexp.MustCompile("(BYTE|WORD|DWORD|QWORD|char)\\s*(\\w+)\\[?([a-zA-Z0-9_\\+]*)\\]?")
}

type Server struct {
	Common.Server
}

func (s *Server) searchMessage(message string, messageFiles []string) map[string]interface{} {

	// form 1
	regBegin := regexp.MustCompile("(MSG_BEGIN_?\\d*)\\([a-zA-Z0-9_]+\\s*,\\s*(\\d+)\\s*,\\s*" + message)

	// form 2
	form2paramNumReg := regexp.MustCompile("const BYTE [a-zA-Z_]+\\s*=\\s*(\\d+)")
	form2structReg := regexp.MustCompile("struct\\s*" + message + "\\s*:\\s*public\\s*([a-zA-Z_\\d]+)")
	form2endReg := regexp.MustCompile("\\s*};")

	response := make(map[string]interface{})
	var params []interface{}

	isForm1 := false
	isForm2 := false

	for _, v := range messageFiles {
		f, err := os.Open(s.Path + "\\" + v)

		if err != nil {
			log.Println(err)
			continue
		}

		r := bufio.NewReader(f)
		var lastLine string

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				break
			}

			// form1
			if !isForm2 {
				ret := regBegin.FindStringSubmatch(line)
				if len(ret) != 0 {
					isForm1 = true
					log.Println("find match:", ret)

					if len(ret) == 3 {
						log.Println("response set")
						response["cmdFlag"] = ret[1] // MSG_BEGIN
						response["paraNumber"] = ret[2]
					}
				}
			}

			// form2
			if !isForm1 {
				ret := form2structReg.FindStringSubmatch(line)
				if len(ret) != 0 {
					log.Println("form2 find match:", ret)
					isForm2 = true
					response["cmdFlag"] = ret[1]
					ret = form2paramNumReg.FindStringSubmatch(lastLine)

					if len(ret) != 0 {
						response["paraNumber"] = ret[1]
						log.Println("form2 find paramcount match:", ret)
					}
				}
			}

			if isForm1 || isForm2 {
				ret := regParam.FindStringSubmatch(line)
				if len(ret) == 4 {
					// ret[1]:paramName ret[2]:type ret[3]:array len
					log.Println("find param match:", ret[1], ret[2], ret[3])

					params = append(params, []string{ret[1], ret[2], ret[3]})
					response["params"] = params
				}
				// form1
				if isForm1 && strings.Contains(line, "MSG_END") {
					break
				}
				// form2
				if isForm2 {
					ret = form2endReg.FindStringSubmatch(line)
					if len(ret) != 0 {
						break
					}
				}
			}
			lastLine = line
		}
		f.Close()

		if isForm1 || isForm2 {
			break
		}
	}

	return response
}
