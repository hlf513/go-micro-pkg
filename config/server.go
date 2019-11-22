package config

type Server struct {
	Name string
	Env  string
	Mode string
}

var server Server

func GetServer() Server {
	return server
}

func SetServer(s Server) {
	server = s
}
