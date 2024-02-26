package main

import (
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

func main() {
	providerFile := flag.String("p", "", "provider file name")
	envFilePath := flag.String("e", "", "env file name")
	flag.Parse()

	getPath, err := filepath.Abs(*providerFile)
	assertErr(err, "failed to build file path")

	logrus.Println("haha", getPath, *providerFile)

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

	cfg, err := loader.ReadFromFile(env, *envFilePath)
	assertErr(err, "failed to load config")

	logrus.Printf("config %v\n", cfg)
}
