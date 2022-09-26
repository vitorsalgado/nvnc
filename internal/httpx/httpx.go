package httpx

import (
	"net/http"
	"time"
)

type Conf struct {
	Timeout time.Duration
}

func Client(conf *Conf) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{Timeout: conf.Timeout, Transport: t}
}
