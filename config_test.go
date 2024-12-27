package diaper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthzConfig struct {
	ClientID     string
	ClientSecret string
	Domain       string
	Audience     string
	CallbackURL  string
}

type AppConfig struct {
	Authz AuthzConfig

	DatabaseURL string
	Port        int
}

func SetupConfigTestEnv(envMap ConfigMap) {
	for key, value := range envMap {
		os.Setenv(key, value.(string))
	}
}

func TearDownConfigTestSuite(envMap ConfigMap) {
	for key := range envMap {
		os.Unsetenv(key)
	}
}

func buildTestConfigStruct(configMap ConfigMap) *AppConfig {
	return &AppConfig{
		Authz: AuthzConfig{
			ClientID:     configMap.MustGet("authz_client_id").(string),
			ClientSecret: configMap.MustGet("authz_client_secret").(string),
			Domain:       configMap.MustGet("authz_domain").(string),
			CallbackURL:  configMap.MustGet("authz_callback_url").(string),
			Audience:     configMap.MustGet("authz_audience").(string),
		},
		DatabaseURL: configMap.MustGet("database_url").(string),
		Port:        configMap.MustGetInt("port"),
	}
}

func Test_ReadConfigFromFile(t *testing.T) {
	t.Run("success when required fields are present", func(t *testing.T) {
		envMap := map[string]interface{}{
			"AUTHZ_CLIENT_ID":     "clientid",
			"AUTHZ_CLIENT_SECRET": "clientsecret",
			"AUTHZ_AUDIENCE":      "audience",
			"AUTHZ_CALLBACK_URL":  "callback_url",
			"DATABASE_URL":        "db_url",
			"AUTHZ_DOMAIN":        "localhost",
		}

		// Set the env values in ENV with
		// os.Setenv()
		// after test complete, unset them.
		SetupConfigTestEnv(envMap)
		defer TearDownConfigTestSuite(envMap)

		dc := DiaperConfig{
			Providers:      Providers{EnvProvider{}},
			DefaultEnvFile: "app.env",
		}

		cfgMap, err := dc.ReadFromFile("test", "./examples/")
		require.NoError(t, err)

		cfg := buildTestConfigStruct(cfgMap)

		assert.Equal(t, "clientid", cfg.Authz.ClientID)
		assert.Equal(t, 9090, cfgMap.MustGetInt("port"))

		assert.NotEqual(t, fmt.Sprintf("%v", cfg.Port), os.Getenv("port"))
	})

	t.Run("set value in OS env from file, if not present", func(t *testing.T) {
		dc := DiaperConfig{
			Providers:      Providers{EnvProvider{}},
			DefaultEnvFile: "app.env",
			SetMissingEnv:  true,
		}

		cfgMap, err := dc.ReadFromFile("test", "./examples/")
		require.NoError(t, err)

		defer func() {
			os.Unsetenv("port")
		}()

		assert.Equal(t, fmt.Sprintf("%v", cfgMap.MustGetInt("port")), os.Getenv("port"))
		assert.Equal(t, cfgMap.MustGetString("authz_client_id"), "")

	})
}
