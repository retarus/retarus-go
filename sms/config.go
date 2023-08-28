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

func NewConfig(user string, password string, region common.Region) Config {
	if user == "" || password == "" {
		log.Fatal("username or password is empty or not set correctly")
		panic("")
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

func NewConfigFromEnv(region common.Region) Config {
	username := os.Getenv("retarus_username")
	password := os.Getenv("retarus_password")
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
