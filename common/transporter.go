package common

import (
	"bytes"
	"net/http"
	"time"
)

type Transporter struct {
	HTTPClient http.Client
}

func NewTransporter(timeout int) Transporter {
	return Transporter{
		http.Client{Timeout: 5 * time.Second},
	}
}

type kvParams struct {
	key, value string
}

func (t *Transporter) DoDatacenterFetch(servers []string, username string, password string, body []byte, resource string, method string, params ...kvParams) []*http.Response {
	ch := make(chan *http.Response)
	// Start two goroutines to fetch the URLs
	for _, baseUrl := range servers {
		go fetch(baseUrl+resource, body, method, username, password, ch)
	}

	responses := []*http.Response{}
	// Wait for results from both goroutines
	for i := 0; i < len(servers); i++ {
		x := <-ch
		responses = append(responses, x)
	}
	return responses
}

func fetch(uri string, body []byte, method string, username string, password string, ch chan<- *http.Response, params ...kvParams) {
	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if len(params) > 0 {
		q := req.URL.Query()
		for _, p := range params {
			q.Add(p.key, p.value)
		}
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		// bsider how you want to handle this error.
		// For now, let's just return without sending anything to the channel.
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if username != "" {
		req.SetBasicAuth(username, password)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		// Handle this error as well. Perhaps you might want to send an error through the channel?
		return
	}
	ch <- res // Send the response via the channel
}
