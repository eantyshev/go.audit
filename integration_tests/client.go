package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const TIMEOUT = 10 * time.Second

type Event struct {
	Id        string                 `json:"id"`
	Type      string                 `json:"type" binding:"required"`
	Consumer  string                 `json:"consumer" binding:"required"`
	CreatedAt time.Time              `json:"created_at"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

type QueryParams struct {
	Type        *string    `json:"type,omitempty"`
	Consumer    *string    `json:"consumer,omitempty"`
	CreatedFrom *time.Time `json:"created_from,omitempty"`
	CreatedTo   *time.Time `json:"created_to,omitempty"`
}

type BaseResponse struct {
	Error  string  `json:"error,omitempty"`
	Events []Event `json:"events,omitempty"`
}

type AuditApiClient struct {
	*http.Client
	baseURL, apiKey string
}

func NewAuditApiClient() *AuditApiClient {
	c := &http.Client{
		Timeout: TIMEOUT,
	}
	baseURL := os.Getenv("AUDIT_API_URL")
	apiKey := os.Getenv("AUDIT_API_KEY")

	client := &AuditApiClient{
		Client:  c,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
	client.waitHttp()

	return client
}

func (c *AuditApiClient) waitHttp() {
	deadline := time.Now().Add(c.Timeout)
	for deadline.After(time.Now()) {
		_, err := c.Get(c.baseURL)
		if err == nil {
			return
		}

		log.Println("waitHttp:", err)
		time.Sleep(time.Second)
	}
}

func (c *AuditApiClient) requestDo(method, url string, req interface{}) (int, *BaseResponse) {
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("method: %v, url: %v, req: \n", method, url)
	json.NewEncoder(os.Stdout).Encode(req)
	fmt.Println()

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Api-Key", c.apiKey)

	resp, err := c.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Printf("%d ", resp.StatusCode)

	tee := io.TeeReader(resp.Body, os.Stdout)
	responseData := new(BaseResponse)
	err = json.NewDecoder(tee).Decode(responseData)
	fmt.Println()

	if err != nil {
		panic(err)
	}

	return resp.StatusCode, responseData
}

func (c *AuditApiClient) AddEvent(event Event) (int, *BaseResponse) {
	return c.requestDo(
		"POST",
		c.baseURL+"/v1/event",
		event,
	)
}

func (c *AuditApiClient) ListEvents(params QueryParams) (int, *BaseResponse) {
	return c.requestDo(
		"GET",
		c.baseURL+"/v1/events",
		params,
	)
}
