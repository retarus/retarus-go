package fax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/retarus/retarus-go/common"
)

// Client is responsible for sending fax requests and handling transportation for a fax service.
// It needs to be properly configured to interact with the service.
// Note: To create a new instance of FaxClient, use the NewClient function.
type Client struct {
	// Config holds the configuration settings for the fax client.
	Config Config

	// Transporter is responsible for the actual HTTP requests and responses.
	Transporter common.Transporter
}

// NewClient creates and returns a new FaxClient instance.
// The client is configured based on the provided Config object.
// The Transporter is initialized with a default timeout of 5 seconds.
//
// This is the preferred way to create a new FaxClient instance.
//
// Parameters:
//   - config: A Config object containing settings like API keys, base URLs, etc.
//
// Returns:
//   - A configured FaxClient object ready to send requests to the fax service.
func NewClient(config Config) Client {
	return Client{
		Config:      config,
		Transporter: common.NewTransporter(5), // Initialize transporter with a timeout of 5 seconds
	}
}

// Send sends a fax job to the specified numbers in the job.
func (c *Client) Send(job Job) (jobID string, err error) {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return "", err
	}
	u, err := url.JoinPath(string(c.Config.Region.HAAddr), "/", c.Config.CustomerNumber, "/fax")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(jobBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.Config.User, c.Config.Password)
	resp, err := c.Transporter.HTTPClient.Do(req)
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

	responses := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, bulkBytes, c.Config.CustomerNumber+"/fax/reports", http.MethodPost)
	var allReports []Report

	for _, resp := range responses {
		defer resp.Body.Close()
		if resp.StatusCode == 200 || resp.StatusCode == 404 {
			var faxReports struct {
				Reports []Report `json:"reports"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&faxReports); err != nil {
				return nil, err
			}

			allReports = append(allReports, faxReports.Reports...)
		} else {
			if err := statusToError(resp.StatusCode, resp.Body); err != nil {
				return nil, err
			}
		}
	}

	return allReports, nil
}

// Takes in a array of job ids and deletes all job reports which can befoudn with the gievn job ids, the response will contain a boolean which verifys if the reprot was delete or not.
func (c *Client) DeleteBulkReports(jobIDs []string) ([]DeleteReport, error) {
	bulkreq := bulkReportRequest{
		Action: "DELETE",
		JobIDs: jobIDs,
	}

	bulkBytes, err := json.Marshal(bulkreq)
	if err != nil {
		return nil, err
	}

	responses := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, bulkBytes, c.Config.CustomerNumber+"/fax/reports", http.MethodPost)
	var allDeletedReports []DeleteReport

	for _, resp := range responses {
		defer resp.Body.Close()
		if resp.StatusCode == 200 || resp.StatusCode == 404 {
			var deleteReport struct {
				Reports []DeleteReport `json:"reports,omitempty"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&deleteReport); err != nil {
				return nil, err
			}
			if deleteReport.Reports[0].Reason == "NOT_FOUND" {
				continue
			}
			allDeletedReports = append(allDeletedReports, deleteReport.Reports...)
		} else {
			if err := statusToError(resp.StatusCode, resp.Body); err != nil {
				return nil, err
			}
		}
	}

	return allDeletedReports, nil
}

// DeleteReports deletes up to 1000 status reports for completed fax jobs for the current account, starting from the
// oldest ones. It returns the jobIds of deleted job reports.
func (c *Client) DeleteReports() ([]DeleteReport, error) {
	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, c.Config.CustomerNumber+"/fax/reports", http.MethodDelete)

	type deleteJobResponse struct {
		Reports []DeleteReport `json:"reports,omitempty"`
	}
	var delReportResp deleteJobResponse

	var deletedReports []DeleteReport

	for _, x := range resp {

		defer x.Body.Close()
		if x.StatusCode == 200 || x.StatusCode == 404 {
			defer x.Body.Close()
			if err := json.NewDecoder(x.Body).Decode(&delReportResp); err != nil {
				return nil, err
			}

			deletedReports = append(deletedReports, delReportResp.Reports...)
		} else {
			if err := statusToError(x.StatusCode, x.Body); err != nil {
				return nil, err
			}
		}
	}
	return delReportResp.Reports, nil
}

// DeleteReport deletes a Report for the given jobID.
func (c *Client) DeleteReport(jobID string) (*DeleteReport, error) {
	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, c.Config.CustomerNumber+"/fax/reports/"+jobID, http.MethodDelete)
	var deleteReport DeleteReport
	for _, x := range resp {
		type reports struct {
			// Reports (required)
			Reports []Report `json:"reports"`
		}
		defer x.Body.Close()
		if x.StatusCode == 404 {
			continue
		}
		if x.StatusCode == 200 {
			defer x.Body.Close()
			if err := json.NewDecoder(x.Body).Decode(&deleteReport); err != nil {
				return nil, err
			}
			return &deleteReport, nil
		}
	}

	return &deleteReport, nil
}

// GetReport gets a Report for the given jobID, GetReport will not delete it
// remotely, use DeleteReport after GetReport.
func (c *Client) GetReport(jobID string) (*Report, error) {
	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, "/"+c.Config.CustomerNumber+"/fax/reports/"+jobID, http.MethodGet)
	var faxReport Report
	for _, x := range resp {
		defer x.Body.Close()
		if x.StatusCode == 404 {
			continue
		}
		if x.StatusCode == 200 {
			defer x.Body.Close()
			if err := json.NewDecoder(x.Body).Decode(&faxReport); err != nil {
				return nil, err
			}
			return &faxReport, nil
		}
	}

	return &faxReport, nil
}

// GetReports fetches available status reports for this account
// Status reports are available for up to 30 days or until deleted.
// Important: The results are limited to the oldes 1000 entries. It is recommended to delete
// the status reports after fetching them in order to retrieve the following ones.
func (c *Client) GetReports() ([]Report, error) {
	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, c.Config.CustomerNumber+"/fax/reports", http.MethodGet)

	var faxReports []Report

	for _, x := range resp {
		type reports struct {
			// Reports (required)
			Reports []Report `json:"reports"`
		}
		defer x.Body.Close()
		if x.StatusCode == 200 || x.StatusCode == 404 {
			var faxReport reports
			defer x.Body.Close()
			if err := json.NewDecoder(x.Body).Decode(&faxReport); err != nil {
				return nil, err
			}

			faxReports = append(faxReports, faxReport.Reports...)
		} else {
			if err := statusToError(x.StatusCode, x.Body); err != nil {
				return nil, err
			}
		}
	}

	return faxReports, nil
}
