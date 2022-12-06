package sms

import (
	"time"
)

// Report is a JobReport specified in 4.1
type Report struct {
	// JobID (required)
	JobID string `json:"jobId"`
	// Src (required)
	Src string `json:"src"`
	// Encoding (required)
	Encoding string `json:"encoding"`
	// Billcode (required)
	Billcode string `json:"billcode"`
	// StatusRequest (required)
	StatusRequested bool `json:"statusRequested"`
	// Flash (required)
	Flash bool `json:"flash"`
	// ValidityMin (required)
	ValidityMin int `json:"validityMin"`
	// CustomerRef (required)
	CustomerRef string `json:"customerRef"`
	// QOS (required)
	QOS string `json:"qos"`
	// ReceiptTS (required)
	ReceiptTS time.Time `json:"receiptTs"`
	// FinishedTs (optional)
	FinishedTS time.Time `json:"finishedTs,omitempty"`
	// RecipientIDs (required)
	RecipientIDs []string `json:"recipientIds"`
}
