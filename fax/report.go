package fax

import (
	"time"
)

type Report struct {
	// JobID (required)
	JobID string `json:"jobId"`
	// RecipientStatus (required)
	RecipientStatus []RecipientStatus `json:"recipientStatus"`
	// Pages (required)
	Pages int `json:"pages"`
	// Reference (optional)
	Reference Reference `json:"reference,omitempty"`
}

type RecipientStatus struct {
	// Number (required)  the fax recipient’s primary number (international format, e.g., +49891234678).
	Number string `json:"number"`
	// Status (required)
	Status string `json:"status"`
	// Reason (required) Explanation of the status.
	Reason string    `json:"reason"`
	SentTS time.Time `json:"sentTs"`
	// DurationInSecs (required) Duration of the fax transmission until received by the fax recipient.
	DurationInSecs int    `json:"durationInSecs"`
	SentToNumber   string `json:"sentToNumber"`
	RemoteCsid     string `json:"remoteCsid"`
}

type bulkReportRequest struct {
	// Action (required) defines the action to be performed on all jobs whose Job ID
	// is provided in the jobIds list
	Action string `json:"action"`
	// JobIDs (required) List of Job IDs to be processed in bulk
	JobIDs []string `json:"jobIds"`
}

type DeleteReport struct {
	// JobID (required) is the Job ID. string
	JobID string `json:"jobId"`
	// Deleted (required) Returns true if the job was successfully deleted, false
	// otherwise. If absent in the response it means the job was
	// successfully deleted.
	Deleted bool `json:"deleted"`
	// Reason (optional): Missing if deletion was successful, otherwise one of the
	// following reason messages is returned:
	// • NOT_FOUND: No report exists for the given job id.
	// • INTERNAL_ERROR: Unspecified server-side error.
	Reason string `json:"reason,omitempty"`
}
