package common

import (
	"errors"
	"fmt"
	"testing"
)

func TestDetermineServiceRegionWithFax(t *testing.T) {
	region := Europe
	service := "fax"
	res, err := DetermineServiceRegion(region, service)
	if err != nil {
		t.Error(err)
	}
	if res.Region != Europe {
		t.Errorf("Wrong RegionURI was returned by DetermineServiceRegions")
	}
}
func TestDetermineServiceRegionWithSms(t *testing.T) {
	region := Europe
	service := "sms"
	res, err := DetermineServiceRegion(region, service)
	if err != nil {
		t.Error(err)
	}
	if res.Region != Europe {
		t.Errorf("Wrong RegionURI was returned by DetermineServiceRegions")
	}
}

// Checks if the service has no such region as request, what will happen.
func TestDetermineServiceRegionWithInvalidRegion(t *testing.T) {
	region := Switzerland
	service := "sms"
	res, err := DetermineServiceRegion(region, service)
	if err == nil {
		if errors.Is(err, errors.New("no RegionURI for this service found")) {
			return
		} else {
			fmt.Println(res)
			t.Errorf("the functions should fail, but it didn't")
		}
	}

}

// Checks what will happen if by accident a wrong "service" is set by the developer or such.
func TestDetermineServiceRegionWithUnknownService(t *testing.T) {
	region := Switzerland
	service := "sms"
	_, err := DetermineServiceRegion(region, service)
	if err == nil {
		if err == errors.New("No RegionURI for this service found.") {
			return
		} else {
			t.Errorf("the functions should fail, but it didn't")
		}

	}
}
