package Common

const (
	ProxyIp    = "127.0.0.1"
	ProxyPort  = ":7000"
	WorkerPort = ":7001"
)

type Server struct {
	Ip        string   `json:"ip"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Users     []string `json:"users"`
	CheckTime int64
}

func (s *Server) Info() string {
	ret := "Name:" + s.Name + ",Ip:" + s.Ip + ",Path:" + s.Path + ",Users:"
	for _, user := range s.Users {
		ret = ret + user + "|"
	}
	return ret
}
