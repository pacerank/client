package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

type Api struct {
	url             string
	headerTransport *addHeaderTransport
	client          *http.Client
}

// This is the default reply structure in all requests
// against the digest service
type DefaultReplyStructure struct {
	CorrelationId string          `json:"correlation_id"`
	Status        int             `json:"status"`
	Service       string          `json:"service"`
	Content       json.RawMessage `json:"content"`
}

type addHeaderTransport struct {
	T     http.RoundTripper
	token *string
}

func (adt *addHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.token != nil {
		req.Header.Add("X-API-KEY", *adt.token)
	}

	bi, _ := debug.ReadBuildInfo()
	req.Header.Add("User-Agent", fmt.Sprintf("PacerankRuntime/%s (%s; %s)", bi.Main.Version, runtime.GOOS, runtime.GOARCH))

	return adt.T.RoundTrip(req)
}

func New(url string) *Api {
	api := &Api{
		url: url,
		headerTransport: &addHeaderTransport{
			T: http.DefaultTransport,
		},
	}

	api.client = &http.Client{
		Transport:     api.headerTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 10,
	}

	return api
}

func (api *Api) AddAuthorizationToken(token string) {
	api.headerTransport.token = &token
}
