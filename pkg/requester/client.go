package requester

import (
	"crypto/tls"
	"net/http"
	"time"
)

type (
	Client interface {
		Get(url string) (*http.Response, error)
	}

	requester struct {
		client *http.Client
	}
)

func NewRequester(timeout int) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Millisecond}
	return &requester{client: client}
}

func (r *requester) Get(url string) (*http.Response, error) {
	return r.client.Get(url)
}
