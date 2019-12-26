package server

import "github.com/hlf513/go-micro-pkg/config"

type Server struct {
	Name  string
	Env   string
}

// server 初始化
var server = &Server{}

// GetServerConf 读取配置
func GetServerConf() (*Server, error) {
	if err := config.GetConfigurator().Get([]string{"server"}, server); err != nil {
		return server, err
	}
	return server, nil
}
