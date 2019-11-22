package config

type Etcd struct {
	Host string
}

var etcd Etcd

func GetEtcd() Etcd {
	return etcd
}

func SetEtcd(conf Etcd) {
	etcd = conf
}
