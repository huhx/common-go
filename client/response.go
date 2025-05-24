package client

import (
	"encoding/xml"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type Response struct {
	Status     string
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (r Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode <= 299
}

func (r Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode <= 499
}

func (r Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode <= 599
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

func XmlBody[R any](response *Response) (*R, error) {
	var result R
	if err := xml.Unmarshal(response.Body, &result); err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return nil, err
	}
	return &result, nil
}
