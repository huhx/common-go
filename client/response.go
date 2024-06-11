package client

import (
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Header     http.Header
	Body       []byte
}

func (r Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode <= 299
}

func (r Response) StringBody() string {
	return string(r.Body)
}

func JsonBody[R any](response *Response) (*R, error) {
	var result R
	if err := json.Unmarshal(response.Body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	return &result, nil
}
