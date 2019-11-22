package config

type Sentry struct {
	Dns string `json:"dns"`
}

var sentry Sentry

func SetSentry(s Sentry) {
	sentry = s
}

func GetSentry() Sentry {
	return sentry
}
