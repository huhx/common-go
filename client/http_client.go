package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	http.Client
}

func NewHttpClient() HttpClient {
	return HttpClient{Client: http.Client{}}
}

func NewTLSHttpClient() HttpClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return HttpClient{Client: http.Client{Transport: transport}}
}

func (c *HttpClient) Get(baseUrl string, params map[string]string, header http.Header) (*Response, error) {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	fullURL := baseUrl + "?" + values.Encode()
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header = header

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
	}, nil
}

func (c *HttpClient) Post(url string, payload interface{}, header http.Header) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header = header

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
	}, nil
}
