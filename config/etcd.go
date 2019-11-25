package config

import "time"

type Etcd struct {
	Host string
	Timeout time.Duration
}

var etcd Etcd

func GetEtcd() Etcd {
	return etcd
}

func SetEtcd(conf Etcd) {
	etcd = conf
}
