package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const limitReadSize = 10 * 1024 * 1024

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
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Printf("Error Parse url: %v", err)
		return nil, err
	}
	values := u.Query()
	for key, value := range params {
		values.Set(key, value)
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	if header != nil {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return c.execute(req)
}

func (c *HttpClient) Post(baseUrl string, payload interface{}, header http.Header) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	if header != nil {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return c.execute(req)
}

func (c *HttpClient) Put(urlPath string, payload interface{}, header http.Header) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}

	req, err := http.NewRequest("PUT", urlPath, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	if header != nil {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return c.execute(req)
}

func (c *HttpClient) Patch(urlPath string, payload interface{}, header http.Header) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}
	req, err := http.NewRequest("PATCH", urlPath, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	if header != nil {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return c.execute(req)
}

func (c *HttpClient) Delete(urlPath string, header http.Header) (*Response, error) {
	req, err := http.NewRequest("DELETE", urlPath, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	if header != nil {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return c.execute(req)
}

func (c *HttpClient) execute(req *http.Request) (*Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	limitReader := io.LimitReader(resp.Body, limitReadSize)
	body, err := io.ReadAll(limitReader)
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
