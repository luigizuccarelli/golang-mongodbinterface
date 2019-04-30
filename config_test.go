package main

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// this test should pass
	_, err := Init("tests/config.json")
	if err != nil {
		t.Errorf("Config Init failed %v ", err)
	}

	// this should fail
	_, err = Init("tests/config-no-file.json")
	if err == nil {
		t.Errorf("Config Init should have failed file not found")
	}

	// this should fail
	_, err = Init("tests/config-errors.json")
	if err == nil {
		t.Errorf("Config Init should have failed file not found")
	}

	// this should fail
	_, err = Init("tests/config-parse-error.json")
	if err == nil {
		t.Errorf("Config Init should have failed file not found")
	}

}

func TestFileNotFound(t *testing.T) {
	_, err := ReadFile("tests/config-no-file.json")
	if err == nil {
		t.Errorf("Config Readfile should fail with file not found")
	}
}

func TestParseJson(t *testing.T) {
	var dummy []byte
	_, err := ParseJson(dummy)
	if err == nil {
		t.Errorf("Config ParseJson should fail with unexpected end of JSON string")
	}
}

func TestMissingFields(t *testing.T) {

	var redis = Redis{Host: "redis", Port: "6397"}
	var config = Config{Level: "", Port: "9000", RedisDB: redis}

	_, err := ValidateJson(config)
	if err == nil {
		t.Errorf("Config ValidateJson should have failed")
	}

	config = Config{Level: "Level", Port: "", RedisDB: redis}
	_, err = ValidateJson(config)
	if err == nil {
		t.Errorf("Config ValidateJson should have failed")
	}

	redis = Redis{Host: "", Port: "6397"}
	config = Config{Level: "Level", Port: "9000", RedisDB: redis}
	_, err = ValidateJson(config)
	if err == nil {
		t.Errorf("Config ValidateJson should have failed")
	}

	redis = Redis{Host: "redis", Port: ""}
	config = Config{Level: "Level", Port: "9000", RedisDB: redis}
	_, err = ValidateJson(config)
	if err == nil {
		t.Errorf("Config ValidateJson should have failed")
	}

}

func TestEnvars(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("SERVER_PORT", "9999")
	os.Setenv("REDIS_HOST", "new-redis-host")
	os.Setenv("REDIS_PORT", "8888")
	os.Setenv("MW_URL", "middlewareUrl")
	os.Setenv("TOKEN", "new-token")
	redis := Redis{Host: "redis", Port: "1234"}
	config := Config{Level: "debug", Port: "9000", RedisDB: redis}
	cfg, err := ValidateJson(config)
	if err != nil {
		t.Errorf("Config ValidateJson should not fail")
	}

	// we need to assert these
	if cfg.Level != "info" {
		t.Error("Config ValidationJson should have log level set to trace")
	}

}
