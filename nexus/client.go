package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"moul.io/http2curl"
	"net/http"
)

// CreateProxy creates a raw proxy in Nexus
func (c *NexusClient) CreateProxy(proxy RawProxy) (err error) {
	url := fmt.Sprintf("%s/service/extdirect", c.URL)

	uiPayload := UIPayload{
		Action: "coreui_Repository",
		Data: []UIData{{
			Attributes: UIDataAttributes{
				Cleanup:       proxy.Cleanup,
				HTTPClient:    proxy.HTTPClient,
				NegativeCache: proxy.NegativeCache,
				Proxy:         proxy.Proxy,
				Storage:       proxy.Storage,
			},
			RawProxy: proxy,
		}},
		Method: "create",
		Type:   "rpc",
	}
	payload := c.getPayload(uiPayload)

	var request *http.Request
	if request, err = c.newPostRequest(url, payload); err == nil {
		client := http.Client{}
		_, err = client.Do(request)
	}
	return
}

// ListRepositories returns all the repositories of a Nexus server
func (c *NexusClient) ListRepositories() (repos []Repository, err error) {
	url := fmt.Sprintf("%s/service/rest/beta/repositories", c.URL)

	var request *http.Request
	var response *http.Response
	var data []byte
	if request, err = c.newGetRequest(url, nil); err == nil {
		client := http.Client{}
		if response, err = client.Do(request); err == nil && response.StatusCode == http.StatusOK {
			if data, err = ioutil.ReadAll(response.Body); err == nil {
				fmt.Println(string(data))
				err = json.Unmarshal(data, &repos)
			}
		}
	}
	return
}

func (c *NexusClient) newPostRequest(url string, payload io.Reader) (request *http.Request, err error) {
	if request, err = c.newRequest(http.MethodPost, url, payload); err == nil {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("accept", "application/json")
	}
	return
}

func (c *NexusClient) newGetRequest(url string, payload io.Reader) (request *http.Request, err error) {
	request, err = c.newRequest(http.MethodGet, url, payload)
	return
}

func (c *NexusClient) newRequest(method, url string, payload io.Reader) (request *http.Request, err error) {
	if request, err = http.NewRequest(method, url, payload); err == nil {
		request.SetBasicAuth(c.Username, c.Password)

		if cmd, debugErr := http2curl.GetCurlCommand(request); debugErr == nil {
			fmt.Println(cmd)
		} else {
			fmt.Printf("cannot convert a HTTP request to a curl command, %#v\n", debugErr)
		}
	} else {
		fmt.Println("HTTP request", url)
	}
	return
}

func (c *NexusClient) getPayload(obj interface{}) (payload io.Reader) {
	if data, err := json.Marshal(obj); err == nil {
		payload = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return
}
