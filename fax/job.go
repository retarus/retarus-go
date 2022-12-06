package fax

import (
	"fmt"
	"time"
)

// Job is a Faxjob specified in 4.5. FaxJobRequest.
type Job struct {
	// Reference (optional)
	Reference *Reference `json:"reference,omitempty"`
	// Documents (optional)
	Documents []Document `json:"documents,omitempty"`
	// TransportOptions (optional)
	TransportOptions *TransportOptions `json:"transportOptions,omitempty"`
	// RenderingOptions (optional)
	RenderingOptions *RenderingOptions `json:"renderingOptions,omitempty"`
	// StatusReportOptions (optional)
	StatusReportOptions *StatusReportOptions `json:"statusReportOptions,omitempty"`
	// Meta (optional)
	Meta *Meta `json:"meta,omitempty"`
	// Recipients (required)
	Recipients []Recipient `json:"recipients"`
}

type Reference struct {
	// CustomerDefinedID (optional) Freely-defined ID string (max. 256 characters).
	CustomerDefinedID string `json:"customerDefinedId,omitempty"`
	// BillingCode (optional) Information on the cost center; format is arbitrary (max. 80 characters).
	BillingCode string `json:"billingCode,omitempty"`
	// BillingInfo (optional) Additional data for internal customer accounting (max. 80 characters).
	BillingInfo string `json:"billingInfo,omitempty"`
}

// RecipientProperty specified in 4.16.
// A cover page can be personalized for each individual recipient and then attached to the front of each
// fax document. This object allows specifying a value for each of the keys in the template.
// See Cover Page Personalization for additional information.
type RecipientProperty struct {
	// Key (required) Name of the key.
	Key string `json:"key"`
	// Value (required) The value assigned to the key.
	Value string `json:"value"`
}

// Recipient is specified in 4.7. FaxRecipient
// If a fax can be successfully sent to the recipient fax number, the transmission is completed and
// considered successful. The destination number which received the fax will be indicated in the report
// data for each recipient under RecipientStatus -> sentToNumber.
type Recipient struct {
	// Number (required) is the dialed number (international format, e.g., +12015551000).
	// Example : "0012012051598
	Number string `json:"number"`
	// AlternativeNumbers (optional) is used as an alternative Number.
	AlternativeNumbers []string `json:"alternativeNumbers,omitempty"`
	// Properties (optional) is Personalized data for the cover page.
	Properties []RecipientProperty `json:"properties,omitempty"`
}

// Charset can be used in the Documents struct.
type Charset string

const (
	DEFAULT      Charset = UTF_8
	US_ASCII     Charset = "US-ASCII"
	UTF_8        Charset = "UTF-8"
	UTF_16       Charset = "UTF-16"
	UTF_16BE     Charset = "UTF-16BE"
	UTF_16LE     Charset = "UTF-16LE"
	ISO_8859_1   Charset = "ISO-8859-1"
	WINDOWS_1252 Charset = "Windows-1252"
)

// Document is specified in 4.4 in the API v1 Documentation (DocumentWithData).
type Document struct {
	// Name (required) is the document file name; the file extension is important for
	// determining the file type, e.g., Invoice-2017-01.pdf. Please
	// note: The maximum possible length of a file name is 32
	// characters. Allowed characters are: a-zA-Z0-9-_. , and no
	// whitespaces, slashes, or other special characters are
	// permitted.
	// Example : "test-document-inline-byte-array.txt"
	Name string `json:"name"`
	// Charset (optional) Character encoding of plain text documents (*.txt)
	Charset Charset `json:"charset,omitempty"`
	// Reference (optional) is an URL pointing to the document to be
	// transmitted.
	// Please note: either Reference or Data can be used in a single document, but not both at the same
	Reference string `json:"reference,omitempty"`
	// Data is a base64 string with data, if no reference is provided. If both
	// are provided, the reference data (see above) is used.
	Data string `json:"data,omitempty"`
}

// TransportOptions contains information on the transmission of the fax.
type TransportOptions struct {
	// CsId (optional) is the sender ID the received fax was sent from (max. 20 characters).
	CsID string `json:"csid,omitempty"`
	// IsExpress (optional) Dlag for transmissions sent express
	IsExpress bool `json:"isExpress,omitempty"`
	// IsBlacklistEnabled (optional)  is a flag for the use of the Robinson List (only for numbers in
	// Germany), ECOFAX (for numbers in France), or Retarus own blacklist.
	IsBlacklistEnabled bool `json:"isBlacklistEnabled,omitempty"`
}

// Mode is the overlay mode.
type Mode string

const (
	// ALL_PAGES the overlay is applied to all pages.
	ALL_PAGES Mode = "ALL_PAGES"
	// NO_OVERLAY no overlay is used (returns the same
	// result as if "no overlay" had been specified in the
	// options)
	NO_OVERLAY Mode = "NO_OVERLAY"

	// FIRST_PAGE the overlay is applied only to the first
	// page (if you are using a cover page, it is considered the
	// first page)
	FIRST_PAGE Mode = "FIRST_PAGE"

	// LAST_PAGE the overlay is applied only to the last page
	LAST_PAGE Mode = "LAST_PAGE" //

	// ALL_BUT_FIRST_PAGE the overlay is applied to all
	// pages except for the first (if you are using a cover page,
	// the overlay will be applied to all other pages because the
	// cover page is considered the first page)
	ALL_BUT_FIRST_PAGE Mode = "ALL_BUT_FIRST_PAGE"

	// ALL_BUT_LAST_PAGE the overlay is applied to allpages except the last one
	ALL_BUT_LAST_PAGE Mode = "ALL_BUT_LAST_PAGE"

	// ALL_BUT_FIRST_AND_LAST_PAGE the overlay is applied to all pages except for the first and the last (the
	// cover page is considered the first page if this mode is
	// used)
	ALL_BUT_FIRST_AND_LAST_PAGE Mode = "ALL_BUT_FIRST_AND_LAST_PAGE"

	// FIRST_FILE if the faxed document consists of multiple
	// files, the overlay will only be used on the first file’s pages
	// (the cover page is considered not to belong to any file
	// and does not an overlay in this mode)
	FIRST_FILE Mode = "FIRST_FILE" //
)

// Overlay is the setting for the overlay (e.g., stationery). A template (e.g., with letter header and footer) can be
// applied to all or specific pages in the fax. A template consists of a one-page, black-and-white
// document. In order to install an overlay, the customer transfers a template to Retarus, and the
// template is then saved in Retarus' infrastructure under a mutually agreed upon name
type Overlay struct {
	// Name (required) the template name, without the path information and file extension.
	// Example : "overlay_template1"
	Name string `json:"name"`
	// Mode (required)
	Mode Mode `json:"mode"`
}

type PaperFormat string

const (
	A4     PaperFormat = "A4"
	Letter PaperFormat = "Letter"
)

type Resolution string

const (
	High Resolution = "HIGH"
	Low  Resolution = "LOW"
)

type RenderingOptions struct {
	// PaperFormat (required)
	PaperFormat PaperFormat `json:"paperFormat"`
	// Resolution (optional)
	Resolution Resolution `json:"resolution,omitempty"`
	// CoverpageTemplate (optional) is the name of the cover page’s template; e.g., coverpagedefault.ftl.html
	CoverpageTemplate string `json:"coverpageTemplate,omitempty"`
	// Overlay (optional)
	Overlay *Overlay `json:"overlay,omitempty"`
	// Header (optional) the content of the header, including control characters.
	// Example : "%tz=CEST Testfax: CSID: %C Recipient number: %# Date: %d.%m.%Y %H:%M %z"
	Header string `json:"header,omitempty"`
}

type AttachedFaxImageFormat string

const (
	// TIFF (default) Fax image is attached as TIFF
	TIFF AttachedFaxImageFormat = "TIFF"
	// PDF: Fax image is attached as PDF
	PDF AttachedFaxImageFormat = "PDF"
	// PDF_WITH_OCR: Fax image is ttached as a searchable PDF file. Additional costs may occur
	PDF_WITH_OCR AttachedFaxImageFormat = "PDF_WITH_OCR"
)

type AttachedFaxImageMode string

const (
	// NEVER (default) attach the fax image
	NEVER AttachedFaxImageMode = "NEVER"
	// SUCCESS_ONLY: Only attach the fax image in case of successful transmission
	SUCCESS_ONLY AttachedFaxImageMode = "SUCCESS_ONLY"
	// FAILURE_ONLY: Only attach the fax image in case of failed transmission
	FAILURE_ONLY AttachedFaxImageMode = "FAILURE_ONLY"
	// ALWAYS attach the fax image
	ALWAYS AttachedFaxImageMode = "ALWAYS"
)

// ReportMail: In addition to querying via Webservice, it is possible to request
// notification for each fax job as soon as processing is completed. The status
// information can be sent by either HTTP POST or email. Separate
// email addresses can each be specified for delivery and failed delivery confirmations. If an email
// address is deleted for either type of confirmation, no notification email will be sent for the confirmation
// type that was removed. The report emails' format is specified through a template which is filled out
// with relevant data (Job ID, job status, details on the fax recipients). A default template is available for
// all customers; however, you can install a customized template. Templates must be encoded in UTF-8
// format. In addition, it is possible to specify whether the fax image should be attached to the report or
// not and if so in which format.
type ReportMail struct {
	// SuccessAddress (optional) Email address, to which delivery confirmations notifications
	// should be sent.Example : "john.doe@retarus.com"
	SuccessAddress string `json:"successAddress,omitempty"`
	// FailureAddress (optional) Email address, to which a notification should be sent when
	// errors occur. Example : "jane.doe@retarus.com"
	FailureAddress string `json:"failureAddress,omitempty"`
	// AttachedFaxImageFormat (optional) Determines when the fax image will be attached to
	// the email
	AttachedFaxImageFormat AttachedFaxImageFormat `json:"attachedFaxImageFormat,omitempty"`
	// AttachedFaxImageMode (optinoal) Determines when the fax image will be attached to the email.
	AttachedFaxImageMode AttachedFaxImageMode `json:"attachedFaxImageMode,omitempty"`
}

type HTTPStatusPush struct {
	// TargetURL (required) Push URL. Example : "http://retarus.com/test-path/test-target"
	TargetURL string `json:"targetUrl"`
	// Principal (optional) The principal/user name.
	Principal string `json:"principal,omitempty"`
	// Credentials (optional) The credential for the principal/username.
	Credentials string `json:"credentials,omitempty"`
	// AuthMethod (optional) default NONE
	AuthMethod AuthMethod `json:"authMethod,omitempty"`
}

type AuthMethod string

const (
	// NONE is the default
	NONE        AuthMethod = "NONE"
	HTTP_BASIC  AuthMethod = "HTTP_BASIC"
	HTTP_DIGEST AuthMethod = "HTTP_DIGEST"
	OAUTH2      AuthMethod = "OAUTH2"
)

type ISO8601Time time.Time

func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05.999-0700"))

	return []byte(stamp), nil
}

// StatusReportOptions settings for the status report. Consists of reportPurgeTs and reportMail.
type StatusReportOptions struct {
	// ReportPurgeTS (required) Not currently valid. The date after which the status report is
	// no longer available. In ISO 8601 format Example : "2018-11-03T20:14:37.098+02:00"
	ReportPurgeTS ISO8601Time `json:"reportPurgeTs"`
	// ReportMail (optional)
	ReportMail *ReportMail `json:"reportMail,omitempty"`
	// HTTPStatusPush (optional)
	HTTPStatusPush *HTTPStatusPush `json:"httpStatusPush,omitempty"`
}

// JobValid Contains the valid start/end of a fax job (in ISO 8601 format). If this data is not defined correctly, you
// will receive a Job Expiration error (HTTP status code 400)
type JobValid struct {
	// Start (optional) is the beginning of validity for the job (in ISO 8601 format). Please
	// note: if this time is not defined correctly, you will receive a
	// Job Expiration error. Example values are "Z" for UTC or
	// -05:00 for EST.
	// By default jobs are immediately valid.
	Start string `json:"start,omitempty"`
	// End (optional) of validity for the job (in ISO 8601 format). Please note
	// that also durations are supported; the following values are all
	// valid expiration times:
	// • 2018-10-11T15:50:21.372Z (expiration set to the exact
	// moment specified)
	// • PT80M (Expiration set to now + 80 minutes)
	// By default jobs expire one month after they begin being valid.
	End string `json:"end,omitempty"`
}

// Meta information about the request.
type Meta struct {
	// CustomerReference (required) is an information that the customer can use for internal references.
	CustomerReference string `json:"customerReference"`
	// JobValid (required) contains the valid start/end of a fax job (in ISO
	// 8601 format). If this data is not defined correctly,
	// you will receive a Job Expiration error (HTTP status code 400)
	JobValid JobValid `json:"jobValid,omitempty"`
}
