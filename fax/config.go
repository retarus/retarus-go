package fax

type Config struct {
	User           string
	Password       string
	CustomerNumber string
	Endpoint       Endpoint
}

type Endpoint string

const (
	DE  Endpoint = "https://faxws-ha.de.retarus.com/rest/v1"
	DE1 Endpoint = "https://faxws.de1.retarus.com/rest/v1"
	DE2 Endpoint = "https://faxws.de2.retarus.com/rest/v1"

	CH  Endpoint = "https://faxws-ha.ch.retarus.com/rest/v1"
	CH1 Endpoint = "https://faxws.ch1.retarus.com/rest/v1"

	SG  Endpoint = "https://faxws.sg1.retarus.com/rest/v1"
	SG1 Endpoint = "https://faxws.sg1.retarus.com/rest/v1"

	US  Endpoint = "https://faxws-ha.us.retarus.com/rest/v1"
	US1 Endpoint = "https://faxws.us1.retarus.com/rest/v1"
	US2 Endpoint = "https://faxws.us2.retarus.com/rest/v1"
)
