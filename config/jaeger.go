package config

type Jaeger struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var jaeger Jaeger

func GetJaeger() Jaeger {
	return jaeger
}

func SetJaeger(j Jaeger) {
	jaeger = j
}
