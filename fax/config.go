package fax

import (
	"log"
	"os"

	"github.com/retarus/retarus-go/common"
)

type Config struct {
	User           string
	Password       string
	CustomerNumber string
	Region         *common.RegionURI
}

// NewConfig initializes and returns a Config instance based on the provided parameters.
//
// Parameters:
//   - user: User credential for authentication.
//   - password: Corresponding password for the user.
//   - customerNumber: The unique identifier for the customer.
//   - region: The target service region.
//
// Returns:
//
//	A populated Config object.
func NewConfig(user string, password string, customerNumber string, region common.Region) Config {
	rg, err := common.DetermineServiceRegion(region, "fax")
	if err != nil {
		panic(err)
	}
	return Config{
		User:           user,
		Password:       password,
		CustomerNumber: customerNumber,
		Region:         rg,
	}
}

// NewConfigFromEnv CreateFaxConfig initializes a new FaxConfig using environment variables.
// It fetches 'retarus_fax_username', 'retarus_fax_password', and 'retarus_cuno'
// from the environment to set up and authenticate with the fax SDK.
// Ensure these variables are correctly set in your environment.
// Parameters:
//   - region: The target service region.
//
// Returns:
//
//	A populated Config object.
func NewConfigFromEnv(region common.Region) Config {
	username := os.Getenv("retarus_fax_username")
	password := os.Getenv("retarus_fax_password")
	customerNumber := os.Getenv("retarus_cuno")

	if username == "" || password == "" || customerNumber == "" {
		log.Fatal("One of mandatory env keys are not set, check if following keys set: retarus_fax_username , retarus_fax_password , retarus_cuno")
		panic("")
	}
	rg, err := common.DetermineServiceRegion(region, "fax")
	if err != nil {
		panic(err)
	}

	return Config{
		User:           username,
		Password:       password,
		CustomerNumber: customerNumber,
		Region:         rg,
	}
}
