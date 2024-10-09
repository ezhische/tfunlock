package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Key       string
	Secret    string
	TableName string
	Region    string
	Args      []string
	Endpoint  string
}

func NewConfig() (*Config, error) {
	var lockTable string
	staticAccessKey, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return nil, fmt.Errorf("environment variable AWS_ACCESS_KEY_ID is not set")
	}
	staticSecretKey, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return nil, fmt.Errorf("environment variable AWS_SECRET_ACCESS_KEY is not set")
	}
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		region = "ru-central1"
	}
	//Endpoint
	customEndpoint, ok := os.LookupEnv("CUSTOM_ENDPOINT")
	if !ok {
		return nil, fmt.Errorf("environment variable CUSTOM_ENDPOINT is not set")
	}
	//Таблица с блокировками
	flag.StringVar(&lockTable, "table", "sample_lock_table", "Table name")
	flag.Parse()
	arg := flag.Args()
	if len(arg) == 0 {
		arg = []string{""}
	}
	return &Config{
		Key:       staticAccessKey,
		Secret:    staticSecretKey,
		TableName: lockTable,
		Region:    region,
		Args:      arg,
		Endpoint:  customEndpoint,
	}, nil
}
