package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/retarus/retarus-go/common"
	"net/http"
	"net/url"
)

type Client struct {
	Config      Config
	Transporter common.Transporter
}

func NewClient(config Config) Client {
	return Client{
		Config:      config,
		Transporter: common.NewTransporter(5),
	}
}

// Send sends a sms job to the specified numbers in the job.
func (c *Client) Send(job Job) (jobID string, err error) {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	resp, err := c.doHTTPRequest(jobBytes, "/jobs", http.MethodPost)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return "", err
	}

	type jobResp struct {
		JobID string `json:"jobId"`
	}

	var jobResponse jobResp
	if err := json.NewDecoder(resp.Body).Decode(&jobResponse); err != nil {
		return "", err
	}

	return jobResponse.JobID, nil
}

type kvParams struct {
	key, value string
}

func (c *Client) doHTTPRequest(body []byte, resource string, method string, params ...kvParams) (*http.Response, error) {
	u, err := url.JoinPath(c.Config.Region.HAAddr, resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.Config.User, c.Config.Password)

	if len(params) > 0 {
		q := req.URL.Query()
		for _, p := range params {
			q.Add(p.key, p.value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return c.Transporter.HTTPClient.Do(req)
}

// GetReport gets a Report for the given jobID.
func (c *Client) GetReport(jobID string) (*Report, error) {
	var smsReport Report

	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, "/jobs/"+jobID, http.MethodGet)
	if len(resp) == 0 {
		return nil, errors.New("Error occured during fetch of responses.")
	}
	for x := range resp {

		//fmt.Println("Server: ", resp[x].Request.URL, "Response-Yeah: ", resp[x].StatusCode)
		if resp[x].StatusCode == 404 {
			continue
		}
		defer resp[x].Body.Close()

		if err := statusToError(resp[x].StatusCode, resp[x].Body); err != nil {
			return nil, err
		}

		if err := json.NewDecoder(resp[x].Body).Decode(&smsReport); err != nil {
			return nil, err
		}
	}
	if smsReport.IsZero() == true {
		return nil, errors.New("no reports found, try again later or contact customer service")
	}
	return &smsReport, nil
}
