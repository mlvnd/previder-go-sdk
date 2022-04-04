package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultBaseURL = "https://portal.previder.nl/api/"
	iaasBasePath   = "v2/iaas/"
	jsonEncoding   = "application/json; charset=utf-8"
)

type BaseClient struct {
	httpClient     *http.Client
	clientOptions  *Options
	Task           TaskService
	VirtualMachine VirtualMachineService
	VirtualNetwork VirtualNetworkService
}

type ApiInfo struct {
	Version string `json:"result,omitempty"`
}

type ApiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ApiErrorResponseBody struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Options struct {
	Token   string
	BaseUrl string
}

type Page struct {
	TotalPages       int
	TotalElements    int
	NumberOfElements int
	Size             int
	Number           int
	Content          json.RawMessage
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

//noinspection GoUnusedExportedFunction
func New(options *Options) (*BaseClient, error) {
	if options.Token == "" {
		return nil, fmt.Errorf("missing token")
	}
	if options.BaseUrl == "" {
		options.BaseUrl = defaultBaseURL
	}

	c := &BaseClient{httpClient: http.DefaultClient, clientOptions: options}
	c.Task = &TaskServiceOp{client: c}
	c.VirtualMachine = &VirtualMachineServiceOp{client: c}
	c.VirtualNetwork = &VirtualNetworkServiceOp{client: c}
	return c, nil
}

func (c *BaseClient) Get(url string, responseBody interface{}) error {
	return c.request("GET", url, nil, &responseBody)
}

func (c *BaseClient) Delete(url string, responseBody interface{}) error {
	return c.request("DELETE", url, nil, &responseBody)
}

func (c *BaseClient) Post(url string, requestBody, responseBody interface{}) error {
	return c.request("POST", url, &requestBody, &responseBody)
}

func (c *BaseClient) Put(url string, requestBody, responseBody interface{}) error {
	return c.request("PUT", url, &requestBody, &responseBody)
}

func (c *BaseClient) request(method string, url string, requestBody, responseBody interface{}) error {
	var b *bytes.Buffer
	if requestBody != nil {
		b = bytes.NewBuffer(nil)
		err := json.NewEncoder(b).Encode(requestBody)
		if err != nil {
			return fmt.Errorf("request: %w", err)
		}
	}
	req, err := http.NewRequest(method, c.clientOptions.BaseUrl+url, b)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", jsonEncoding)
	req.Header.Set("X-Auth-Token", c.clientOptions.Token)
	req.Header.Set("Accept", jsonEncoding)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] [Previder API] Error from Previder API received: %s", err.Error())
		return err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiError := new(ApiError)

		var apiErrorResponseBody ApiErrorResponseBody
		temp, err := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(temp, &apiErrorResponseBody)
		if err != nil {
			log.Printf("[ERROR] [Previder API] Could not parse error result: %s", string(temp))
			return err
		}
		log.Printf("[ERROR] [Previder API] Error while executing the request to %s: [%d] :%s",
			apiErrorResponseBody.Path, apiErrorResponseBody.Status, apiErrorResponseBody.Message)
		apiError.Code = resp.StatusCode
		apiError.Message = "[Previder API] " + apiErrorResponseBody.Message
		return apiError
	}

	if responseBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

func (c *BaseClient) ApiInfo() (*ApiInfo, error) {
	apiInfo := new(ApiInfo)
	err := c.Get("", apiInfo)
	return apiInfo, err
}
