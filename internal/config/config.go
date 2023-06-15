package config

import (
	"os"
	"sync"
)

type EnvConfigStruct struct {
	Region           string
	PrimaryTableName string
}

var Config *EnvConfigStruct
var onceConfig sync.Once

func GetConfig() *EnvConfigStruct {
	onceConfig.Do(func() {

		Config = &EnvConfigStruct{
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
