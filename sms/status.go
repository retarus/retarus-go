package sms

// SmsStatus represents the status of an SMS.
// It includes fields like smsId, destination, process status, etc.
type SmsStatus struct {
	smsId         string `json:"smsId"`
	dst           string `json:"dst"`
	processStatus string `json:"processStatus"`
	status        string `json:"status"`
	customerRef   string `json:"customerRef"`
	reason        string `json:"reason"`
	sentTs        string `json:"sentTs"`
	finishedTs    string `json:"finishedTs"`
}

func (s SmsStatus) IsZero() bool {
	return s.smsId == "" &&
		s.dst == "" &&
		s.processStatus == "" &&
		s.status == "" &&
		s.customerRef == "" &&
		s.reason == "" &&
		s.sentTs == "" &&
		s.finishedTs == ""
}
