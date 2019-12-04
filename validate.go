package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// checkEnvars - private function, iterates through each item and checks the required field
func checkEnvar(item string) error {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	logger.Trace(fmt.Sprintf("name %s : required %t", name, required))
	if os.Getenv(name) == "" {
		if required {
			logger.Error(fmt.Sprintf("%s envar is mandatory please set it", name))
			return errors.New(fmt.Sprintf("%s envar is mandatory please set it", name))
		} else {
			logger.Error(fmt.Sprintf("%s envar is empty please set it", name))
		}
	}
	return nil
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
// Each microservice will obviously have a diffefrent envars so change where needed
func ValidateEnvars() error {
	items := []string{
		"LOG_LEVEL,false",
		"SERVER_PORT,true",
		"REDIS_HOST,true",
		"REDIS_PORT,true",
		"REDIS_PASSWORD,true",
		"MONGODB_HOST,true",
		"MONGODB_DATABASE,true",
		"MONGODB_USER,true",
		"MONGODB_PASSWORD,true",
		"VERSION,true",
		"URL,true",
		"PROVIDER_NAME,true",
		"PROVIDER_URL,true",
		"PROVIDER_TOKEN,true",
		"ANALYTICS_URL,true",
	}
	for x, _ := range items {
		if err := checkEnvar(items[x]); err != nil {
			return err
		}
	}
	return nil
}
