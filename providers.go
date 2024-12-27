package diaper

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProviderConfig struct {
	ProviderKeys []string `yaml:"provider"`
}

func LoadProviders(reader io.Reader) Providers {
	pc := ProviderConfig{}

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&pc); err != nil {
		log.Fatal(fmt.Errorf("failed to decode provider config. error %w", err))
	}

	providers := Providers{}

	for _, key := range pc.ProviderKeys {
		switch strings.ToLower(key) {
		case "env":
			providers = append(providers, EnvProvider{})
		default:
			providers = append(providers, NoopProvider{})
		}
	}

	return providers
}

func BuildProviders(providers ...ValueProvider) Providers {
	providers = append(providers, NoopProvider{})
	return providers
}

type ValueProvider interface {
	Deref(interface{}) interface{}
}

type Providers []ValueProvider

func (providers Providers) Deref(value interface{}) interface{} {
	for _, provider := range providers {
		value = provider.Deref(value)
	}

	return value
}

const EnvProviderPrefix = "env://"

type EnvProvider struct{}

func (EnvProvider) Deref(value interface{}) interface{} {
	// At this point the value need to be a string
	// otherwise pass as is

	str, ok := value.(string)
	if !ok {
		return value
	}

	if strings.HasPrefix(str, EnvProviderPrefix) {
		return os.Getenv(
			strings.TrimPrefix(str, EnvProviderPrefix),
		)
	}

	return value
}

type NoopProvider struct{}

func (NoopProvider) Deref(value interface{}) interface{} {
	return value
}

// const SSMProviderPrefix = "ssm://"
// type SSMProvider struct{
//  paramPrefix string
//  awsClient string
//  fetchedValues map[string]string{}
// }

// func NewSSMProvider(awsClient *ssm.Client, paramPrefix string) SSMProvider {
//   fetchedValues := awsClient.GetParameetersByPath(paramPrefix)
//}
