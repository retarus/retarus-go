package sms

import (
	"fmt"
	"time"
)

// Job is a SMS Job Request specified in 4.2.
type Job struct {
	// Messages (required) SMS messages with different texts and recipients. A
	// minimum of one message must be specified in the list.
	// Please note: Requests which result in more than 3.000 SMS
	// (based on several messages or recipients) will be rejected.
	Messages []Message `json:"messages"`
	// Options (optional)
	Options *Options `json:"options,omitempty"`
}

func NewJob(message []Message, options *Options) Job {
	return Job{
		message,
		options,
	}
}

// AddMessage will append a message to the Job.
func (j *Job) AddMessage(m Message) {
	j.Messages = append(j.Messages, m)
}

// AddMessages add an array of Messages to the job, but will override all already set messages
func (j *Job) AddMessages(m []Message) {
	j.Messages = m
}

type Recipient struct {
	// Dst (required) Recipient’s mobile phone number. If a number is specified
	// without a country code, the value specified in the Country
	// code parameter will automatically be applied.
	Dst string `json:"dst"`
	// CustomerRef (optional) Discretionary reference information. If no value is specified,
	// the recipient mobile phone number specified in the dst field
	// will be applied. This information is included in the status
	// report. Please note: Max. 70 characters.
	CustomerRef string `json:"customerRef,omitempty"`
	// BlackoutPeroids (optional)
	// 	• Specifies time periods during which no SMS are
	// delivered (in accordance with ISO-8601). E.g. if a
	// delivery time is inside such a period, the SMS is
	// scheduled to the end of the period.
	// • These periods are based on the ISO 8601 standard.
	// • Periods that are smaller than 1 hour are expanded to 1
	// hour (e.g.,: 17:10 - 17:20 is expanded to 17:10 - 18:10)
	// • If the blackoutPeriods syntax is invalid, the job is rejected
	// with a 400 (Bad Request) status code.
	// • If blackout periods are specified at the Recipient level,
	// only they are used. The blackout periods in the Options
	// are then ignored.
	BlackoutPeriods []time.Time `json:"blackoutPeriods,omitempty"`
}

func NewRecipient(destination string, customerRef string, blackout []time.Time) Recipient {
	return Recipient{
		destination,
		customerRef,
		blackout,
	}
}

type Message struct {
	// Text (required) to send
	Text string `json:"text"`
	// Recipients (required)
	Recipients []Recipient `json:"recipients"`
}

func NewMessage(text string, recipients []Recipient) Message {
	return Message{
		text,
		recipients,
	}
}

type Encoding string

const (
	STANDARD Encoding = "STANDARD"
	UTF16    Encoding = "UTF-16"
)

type InvalidCharacters string

const (
	// 	REFUSE rejects the job with status-code 422
	// (UNPROCESSABLE_ENTITY); see also HTTP-statuscodes below
	REFUSE InvalidCharacters = "REFUSE"

	// REPLACE invalid characters by blank space
	REPLACE InvalidCharacters = "REPLACE"

	// TO_UTF16 send SMS as UTF-16. Please note: each
	// part now has a maximum of only 70 characters (67 if
	// concatenated)
	TO_UTF16 InvalidCharacters = "TO_UTF16"

	// TRANSLITERATE: do a smart replacement
	TRANSLITERATE InvalidCharacters = "TRANSLITERATE"
)

type QOS string

const (
	NORMAL  QOS = "normal"
	EXPRESS QOS = "express"
)

// Options (optional) Global job options.
type Options struct {
	// Src (optional) The Sender ID displayed to recipient. When entering the
	// Sender ID, please note the following technical restrictions:
	// • If the Sender ID is exclusively numerical, a maximum of
	// 20 numbers can be used. A + character is available for
	// optional use.
	// • If alphanumeric characters are used, the maximum
	// length is 11. The following characters can be used: a-z,
	// A-Z, 0-9, !"#$%&'()*+,-./:;<=>?@[\]^_{|}~
	// No other characters are permitted. Java Regular
	// Expressions: (.\\+|)[0-9]{1,20}|[a-zA-Z0-
	// 9\\x20\\p{Punct}]{1,11}
	// Please note:
	// • p{Punct} stands for the following: !"#$%&'()*+,-
	// ./:;<=>?@[\]^_{|}~.
	// Pattern : "(\\+|)[0-9]{1,20}|[a-zA-Z0-
	// 9\\x20\\p{Punct}]{1,11}"
	// Example : "retarus"
	Src string `json:"src,omitempty"`
	// Encoding (optional)
	// The preferred encoding of an SMS. Potential values:
	// • STANDARD
	// • UTF-16
	// When STANDARD is selected, the Webservice checks SMS
	// messages for invalid characters.
	// Please note: the selection of encoding impacts the maximum
	// length of the SMS.
	// Important information regarding SMS coding using GSM-7
	// (standard):
	// • There is a maximum of 160 characters per single-part SMS
	// • A maximum of 153 characters per SMS in a multi-part SMS
	// • Some special characters, such as, e.g., €, require
	// additional storage capacity (2 normal characters)
	// Important information regarding SMS encoding with UTF-16:
	// • A maximum of 70 characters when sending a single-part SMS
	// • A maximum of 67 characters per SMS in a multi-part SMS
	// In some countries, an SMS will not reach the recipient if it
	// contains special characters.
	Encoding Encoding `json:"encoding,omitempty"`
	// Billcode (optional) billing information can be entered here. max 70 chars
	Billcode string `json:"billcode,omitempty"`
	// StatusRequested (optional) Requests a delivery notification.
	StatusRequested bool `json:"statusRequested,omitempty"`
	// Flash (optional) can set here whether you want the SMS sent as a Flash
	// SMS, which is sent directly to the recipient’s mobile phone
	// display without the recipient explicitly having to open it.
	// Please note: An additional fee is required. The Flash SMS
	// option is not supported in all countries
	Flash bool `json:"flash,omitempty"`
	// CustomerRef (optional) reference information that will be included in
	// getJobReport.Please note: max. 70 characters (192 for US-ASCII
	// encoding).
	CustomerRef string `json:"customerRef,omitempty"`
	// ValidityMin (optional) Expresses the validity of an SMS in minutes and thus the
	// length of the delivery attempt when an SMS cannot instantly
	// be delivered after being sent (e.g., if the recipient’s mobile
	// phone is shut off). Supported values are between 5 and 2880
	// minutes. If the specified value is outside of the range of
	// permitted values, either the minimum or maximum value will
	// automatically be selected. Exception: 0, in which case the
	// provider’s default value is used.
	ValidityMin int `json:"validityMin,omitempty"`
	// MaxParts (optional) Sets the maximum amount of SMS in a multi-part
	// (concatenated) SMS. If the SMS message is longer, it is
	// truncated. Permitted values: from 1 to 20. If the specified
	// value is outside of the range of permitted values, either the
	// minimum or maximum value will automatically be selected.
	// The actual length of an SMS is dependent on the selected
	// encoding (see encoding earlier in this section for additional
	// information).
	// Please note: Multi-part SMS are not supported in CDMA
	// networks.
	MaxParts int `json:"maxParts,omitempty"`
	// InvalidCharacters (optional) Defines how invalid characters in SMS text will be handled.
	InvalidCharacters InvalidCharacters `json:"invalidCharacters,omitempty"`
	// QOS (optional) Sets the priority of an SMS job. If you have time-critical SMS
	// messages, we recommend that you select the express
	// option.
	QOS string `json:"qos,omitempty"`
	// JobPeriod (optional) Timestamp setting the transmission time of the SMS Job.
	// In accordance with the ISO-8601 standard e.g., Z can be
	// +02:00.
	JobPeriod *ISO8601Time `json:"jobPeriod,omitempty"`
	// DuplicateDetection (optional)
	// • If enabled, equal requests are rejected with a 409
	// (Conflict) status code
	// • The equality is determined from the REST JSON data.
	// • Currently the time window for duplicate detection is 10
	// minutes. This means that after 10 minutes the same job
	// can be successfully sent once again.
	// • This parameter should normally be enabled if the
	// customer’s system erroneously generates duplicate
	// requests.
	DuplicateDetection bool `json:"duplicateDetection,omitempty"`
	// BlackoutPeriods (optional)
	// • Specifies time periods during which no SMS are
	// delivered; e.g. if a delivery time is inside such a period,
	// the SMS is scheduled to be sent at the end of the period.
	// • These periods are based on the ISO 8601 standard.
	// • Periods that are smaller than 1 hour are expanded to 1
	// hour (e.g.,: 17:10 - 17:20 is expanded to 17:10 - 18:10)
	// • If the blackoutPeriods syntax is invalid, the job is rejected
	// with a 400 (Bad Request) status code.
	// • If blackout periods are specified at the Recipient level,
	// only they are used. The blackout periods in the Options
	// are then ignored.
	// Examples:
	// • 2018-10-25T18:00Z/2018-10-26T07:00Z
	// • 2018-10-25T18:00+01:00/2018-10-26T07:00+01:00
	BlackoutPeriods []string `json:"blackoutPeriods,omitempty"`
}

type ISO8601Time time.Time

func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(layout))

	return []byte(stamp), nil
}

const layout = "2006-01-02T15:04:05.999-0700"
