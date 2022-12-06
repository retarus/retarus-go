package fax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	Config     Config
	HTTPClient http.Client
}

// Send sends a fax job to the specified numbers in the job.
func (c *Client) Send(job Job) (jobID string, err error) {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	resp, err := c.doHTTPRequest(jobBytes, "/fax", http.MethodPost)
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

func (c *Client) doHTTPRequest(body []byte, resource string, method string) (*http.Response, error) {
	u, err := url.JoinPath(string(c.Config.Endpoint), "/", c.Config.CustomerNumber, "/", resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.Config.User, c.Config.Password)

	return c.HTTPClient.Do(req)
}

// GetBulkReports It is possible to perform bulk operations on the status reports through a POST .
// The maximum number of jobs per POST request is set to 1000.
func (c *Client) GetBulkReports(jobIDs []string) ([]Report, error) {
	bulkreq := bulkReportRequest{
		Action: "GET",
		JobIDs: jobIDs,
	}

	bulkBytes, err := json.Marshal(bulkreq)
	if err != nil {
		return nil, err
	}

	resp, err := c.doHTTPRequest(bulkBytes, "/fax/reports", http.MethodPost)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	type reports struct {
		// Reports (required)
		Reports []Report `json:"reports"`
	}

	var faxReports reports
	if err := json.NewDecoder(resp.Body).Decode(&faxReports); err != nil {
		return nil, err
	}

	return faxReports.Reports, nil
}

// DeleteBulkReports It is possible to perform bulk operations on the status reports through a POST .
// The maximum number of jobs per POST request is set to 1000.
func (c *Client) DeleteBulkReports(jobIDs []string) ([]DeleteReport, error) {
	bulkreq := bulkReportRequest{
		Action: "DELETE",
		JobIDs: jobIDs,
	}

	bulkBytes, err := json.Marshal(bulkreq)
	if err != nil {
		return nil, err
	}

	resp, err := c.doHTTPRequest(bulkBytes, "/fax/reports", http.MethodPost)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	type deleteJobResponse struct {
		Reports []DeleteReport `json:"reports,omitempty"`
	}

	var deleteReport deleteJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteReport); err != nil {
		return nil, err
	}

	return deleteReport.Reports, nil
}

// DeleteReports deletes up to 1000 status reports for completed fax jobs for the current account, starting from the
// oldest ones. It returns the jobIds of deleted job reports.
func (c *Client) DeleteReports() ([]DeleteReport, error) {
	resp, err := c.doHTTPRequest([]byte{}, "/fax/reports", http.MethodDelete)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	type deleteJobResponse struct {
		Reports []DeleteReport `json:"reports,omitempty"`
	}

	var delReportResp deleteJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&delReportResp); err != nil {
		return nil, err
	}

	return delReportResp.Reports, nil
}

// DeleteReport deletes a Report for the given jobID.
func (c *Client) DeleteReport(jobID string) (*DeleteReport, error) {
	resp, err := c.doHTTPRequest([]byte{}, "/fax/reports/"+jobID, http.MethodDelete)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	var deleteReport DeleteReport
	if err := json.NewDecoder(resp.Body).Decode(&deleteReport); err != nil {
		return nil, err
	}

	return &deleteReport, nil
}

// GetReport gets a Report for the given jobID, GetReport will not delete it
// remotely, use DeleteReport after GetReport.
func (c *Client) GetReport(jobID string) (*Report, error) {
	resp, err := c.doHTTPRequest([]byte{}, "/fax/reports/"+jobID, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	var faxReport Report
	if err := json.NewDecoder(resp.Body).Decode(&faxReport); err != nil {
		return nil, err
	}

	return &faxReport, nil
}

// GetReports fetches available status reports for this account
// Status reports are available for up to 30 days or until deleted.
// Important: The results are limited to the oldes 1000 entries. It is recommended to delete
// the status reports after fetching them in order to retrieve the following ones.
func (c *Client) GetReports() ([]Report, error) {
	resp, err := c.doHTTPRequest([]byte{}, "/fax/reports", http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return nil, err
	}

	type reports struct {
		// Reports (required)
		Reports []Report `json:"reports"`
	}

	var faxReports reports
	if err := json.NewDecoder(resp.Body).Decode(&faxReports); err != nil {
		return nil, err
	}

	return faxReports.Reports, nil
}
