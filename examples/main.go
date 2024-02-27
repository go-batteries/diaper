package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/go-batteries/diaper"
	"github.com/sirupsen/logrus"
)

func assertErr(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}

func setupEnv() func() {
	os.Setenv("DATABASE_NAME", "postgres://database_url")
	return func() {
		os.Unsetenv("DATABASE_NAME")
	}
}

type AppConfig struct {
	DatabaseName string
	AwsAccessKey string
	Port         int
}

func main() {
	providerFile := flag.String("p", "", "provider file name")
	envFilePath := flag.String("e", "", "env file name")
	flag.Parse()

	getPath, err := filepath.Abs(*providerFile)
	assertErr(err, "failed to build file path")

	reader, err := os.Open(getPath)
	assertErr(err, "failed to read file")

	providers := diaper.LoadProviders(reader)

	loader := diaper.DiaperConfig{
		DefaultEnvFile: "app.env",
		Providers:      providers,
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	defer setupEnv()()

	cfgMap, err := loader.ReadFromFile(env, *envFilePath)
	assertErr(err, "failed to load config")

	cfg := AppConfig{
		DatabaseName: cfgMap.MustGet("database_name").(string),
		Port:         cfgMap.MustGetInt("port"),
		AwsAccessKey: cfgMap.MustGet("aws_access_key_id").(string),
	}

	b, _ := json.MarshalIndent(cfg, " ", "  ")
	logrus.Printf("%s\n", string(b))
}
