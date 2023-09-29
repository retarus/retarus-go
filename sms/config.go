package sms

import (
	"github.com/retarus/retarus-go/common"
	"log"
	"os"
)

type Config struct {
	User     string
	Password string
	Region   *common.RegionURI
}

// NewConfig initializes a Config instance using explicitly passed credentials and region.
// This function will panic and terminate the program if either the 'user' or 'password'
// fields are empty or improperly set.
//
// Parameters:
//   - user: The username required for authentication.
//   - password: The corresponding password for the given username.
//   - region: A common.Region enum specifying the target service region.
//
// Returns:
//
//	A fully initialized Config object.
func NewConfig(user string, password string, region common.Region) Config {
	if user == "" || password == "" {
		log.Fatal("username or password is empty or not set correctly")
	}
	rg, err := common.DetermineServiceRegion(region, "fax")
	if err != nil {
		panic(err)
	}
	return Config{
		User:     user,
		Password: password,
		Region:   rg,
	}
}

// NewConfigFromEnv initializes a Config instance by pulling credentials from environment variables.
// It specifically looks for 'retarus_sms_username' and 'retarus_sms_password' in the environment.
// If these variables are not set or are empty, the function will panic and terminate the program.
//
// Parameters:
//   - region: A common.Region enum specifying the target service region.
//
// Returns:
//
//	A fully initialized Config object populated with credentials from the environment.
func NewConfigFromEnv(region common.Region) Config {
	username := os.Getenv("retarus_sms_username")
	password := os.Getenv("retarus_sms_password")
	if username == "" || password == "" {
		log.Fatal("One of mandatory env keys are not set, check if following keys set: retarus_username, retarus_password")
		panic("")
	}
	rg, err := common.DetermineServiceRegion(region, "sms")
	if err != nil {
		panic(err)
	}
	return Config{
		User:     username,
		Password: password,
		Region:   rg,
	}
}
