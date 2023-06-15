package main

import (
	"os"
	"sync"
)

type ConfigStruct struct {
	Region           string
	PrimaryTableName string
}

var Config *ConfigStruct
var onceConfig sync.Once

func GetConfig() *ConfigStruct {
	onceConfig.Do(func() {

		Config = &ConfigStruct{
			Region:           getEnv("AWS_REGION", "eu-west-1"),
			PrimaryTableName: getEnv("DYNAMODB_TABLE_NAME", ""),
		}
	})
	return Config
}

func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}
