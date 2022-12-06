package sms

type Config struct {
	User     string
	Password string
	Endpoint Endpoint
}

type Endpoint string

const (
	EU  Endpoint = "https://sms4a.eu.retarus.com/rest/v1"
	DE1 Endpoint = "https://sms4a.de1.retarus.com/rest/v1"
	DE2 Endpoint = "https://sms4a.de2.retarus.com/rest/v1"
)
