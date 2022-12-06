package sms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Config     Config
	HTTPClient http.Client
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
	u, err := url.JoinPath(string(c.Config.Endpoint), "/", resource)
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

	return c.HTTPClient.Do(req)
}

// GetReport gets a Report for the given jobID.
func (c *Client) GetReport(jobID string) (*Report, error) {
	resp, err := c.doHTTPRequest([]byte{}, "/jobs/"+jobID, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	var smsReport Report
	if err := json.NewDecoder(resp.Body).Decode(&smsReport); err != nil {
		return nil, err
	}

	return &smsReport, nil
}

type Params struct {
	// JobIDsOnly required
	JobIDsOnly bool
	// Limit (optional) default 100
	// Limits the results list to a specific number of
	// Job IDs (0 < limit <=1000).
	Limit int
	// Offset (optional)
	// 	If the number of results is larger than the
	// limit set for it, with the assistance of the
	// offset you can query more recent results or
	// skip over a specified number of Job IDs.
	Offset int
	// Open (optional)
	// Restricts the results list to Job IDs that are
	// either still open or have already been
	// completed (blank = both conditions).
	Open bool
	// FromTS timestamp in ISO-8601 format
	// (maximum 30 days before toTs).
	FromTS *time.Time
	// ToTS optional
	// To timestamp in ISO-8601 format (must be after fromTs)
	ToTS *time.Time
}

// GetJobIDs gets ids for all jobs matching a given criterion.
func (c *Client) GetJobIDs(params Params) ([]string, error) {
	var kvp []kvParams

	switch params.JobIDsOnly {
	case true:
		kvp = append(kvp, kvParams{"jobIdsOnly", "true"})
	case false:
		kvp = append(kvp, kvParams{"jobIdsOnly", "false"})
	}

	if params.FromTS != nil {
		kvp = append(kvp, kvParams{"fromTs", params.FromTS.Format(layout)})
	}

	if params.ToTS != nil {
		kvp = append(kvp, kvParams{"toTs", params.FromTS.Format(layout)})
	}

	if params.Open {
		kvp = append(kvp, kvParams{"open", "true"})
	}

	if params.Offset > 0 {
		kvp = append(kvp, kvParams{"offset", strconv.Itoa(params.Offset)})
	}

	if params.Limit > 0 {
		kvp = append(kvp, kvParams{"limit", strconv.Itoa(params.Limit)})
	}

	resp, err := c.doHTTPRequest([]byte{}, "/jobs", http.MethodGet, kvp...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	type reports struct {
		// jobId (required)
		JobID string `json:"jobId"`
	}

	var smsReports []reports
	if err := json.NewDecoder(resp.Body).Decode(&smsReports); err != nil {
		return nil, err
	}

	jobIds := make([]string, len(smsReports))
	for i, rep := range smsReports {
		jobIds[i] = rep.JobID
	}

	return jobIds, nil
}
