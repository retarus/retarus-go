package common

import "errors"

type Region string

const (
	Europe      Region = "Europe"
	America     Region = "America"
	Switzerland Region = "Switzerland"
	Singapore   Region = "Singapore"
)

type RegionURI struct {
	Region  Region
	HAAddr  string
	Servers []string
}

func NewRegionURI(region Region, haAddr string, servers []string) RegionURI {
	return RegionURI{
		Region:  region,
		HAAddr:  haAddr,
		Servers: servers,
	}
}

func DetermineServiceRegion(region Region, service string) (*RegionURI, error) {
	fax := []RegionURI{
		NewRegionURI(Europe, "https://faxws-ha.de.retarus.com/rest/v1/", []string{"https://faxws.de2.retarus.com/rest/v1/", "https://faxws.de1.retarus.com/rest/v1/"}),
		NewRegionURI(America, "https://faxws-ha.us.retarus.com/rest/v1/", []string{"https://faxws.us2.retarus.com/rest/v1/", "https://faxws.us1.retarus.com/rest/v1/"}),
		NewRegionURI(Switzerland, "https://faxws-ha.ch.retarus.com/rest/v1/", []string{"https://faxws.ch1.retarus.com/rest/v1/"}),
		NewRegionURI(Singapore, "https://faxws.sg1.retarus.com/rest/v1/", []string{"https://faxws.sg1.retarus.com/rest/v1/"}),
	}
	sms := []RegionURI{
		NewRegionURI(Europe, "https://sms4a.eu.retarus.com/rest/v1", []string{"https://sms4a.de1.retarus.com/rest/v1", "https://sms4a.de2.retarus.com/rest/v1"}),
	}
	if service == "sms" {
		for _, uri := range sms {
			if uri.Region == region {
				return &uri, nil
			}
		}
		return nil, errors.New("unknown region type")
	}
	if service == "fax" {
		for _, uri := range fax {
			if uri.Region == region {
				return &uri, nil
			}
		}
		return nil, errors.New("unknown region type")
	}
	return nil, errors.New("no RegionURI for this service found")
}
