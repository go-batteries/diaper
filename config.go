package diaper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type DiaperConfig struct {
	Providers      Providers
	DefaultEnvFile string
	SetMissingEnv  bool
}

const (
	DefaultEnvFile = ".env"
)

var ErrBuildFilePath = errors.New("failed_to_build_file_path")

func (dc *DiaperConfig) ReadFromFile(env, path string) (ConfigMap, error) {
	env = strings.ToLower(env)

	if dc.DefaultEnvFile == "" {
		dc.DefaultEnvFile = DefaultEnvFile
	}

	// find the proper .env file to use
	// its either ${ENV}.env or dc.DefaultEnvFile

	defaultEnvFilePath, err := filepath.Abs(
		filepath.Join(path, dc.DefaultEnvFile),
	)
	if err != nil {
		log.Print(fmt.Errorf("file path %s, failed to load: %w", dc.DefaultEnvFile, err))
		return nil, err
	}

	envOrrideFile := fmt.Sprintf("%s.env", env) // ENV=test -> test.env

	envFilePath, err := filepath.Abs(
		filepath.Join(path, envOrrideFile),
	)
	if err == nil {
		_, err = os.Stat(envFilePath)
	}

	if err != nil {
		// if override file not found use default
		// logrus.WithError(err).Warnln("failed to find override env file", envOrrideFile)
		// logrus.WithError(err).Debugln("file full path tried ", envFilePath)
		envFilePath = defaultEnvFilePath
	}

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(envFilePath)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("failed to read env file. error %w", err))
	}

	configMap := ConfigMap{}

	if err := viper.Unmarshal(&configMap); err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal config. error %w", err))
	}

	for key, value := range configMap {
		configMap[key] = dc.Providers.Deref(value)

		if !dc.SetMissingEnv {
			continue
		}

		// Set the value in ENV if not present
		configValue := configMap.MustGetString(key)
		if os.Getenv(key) != configValue {
			os.Setenv(key, configValue)
		}
	}

	// logrus.Debugln(configMap)
	return configMap, nil
}

type ConfigMap map[string]interface{}

func (cmap ConfigMap) Get(key string) (interface{}, bool) {
	value, ok := cmap[key]
	if !ok {
		return nil, false
	}

	return value, true
}

func (cmap ConfigMap) GetString(key string) (string, bool) {
	value, ok := cmap[key]
	if !ok {
		return "", false
	}

	strValue, ok := value.(string)
	if !ok {
		return "", false
	}
	return strValue, true
}

func (cmap ConfigMap) GetInt(key string) (int, bool) {
	value, ok := cmap[key]
	if !ok {
		return -1, false
	}

	intVal, ok := value.(int)
	if ok {
		return intVal, ok
	}

	strVal, ok := value.(string)
	if !ok {
		return -1, ok
	}

	var err error

	intVal, err = strconv.Atoi(strVal)
	if err != nil {
		return -1, false
	}

	return intVal, true
}

func (cmap ConfigMap) MustGetInt(key string) int {
	value, ok := cmap.GetInt(key)
	if !ok {
		log.Fatal("value for", key, "cannot be cooerced to string")
	}

	return value
}

func (cmap ConfigMap) MustGet(key string) interface{} {
	value, ok := cmap.Get(key)
	if !ok {
		log.Fatal(key, "not found")
	}

	return value
}

func (cmap ConfigMap) MustGetString(key string) string {
	value, ok := cmap.GetString(key)
	if !ok {
		return fmt.Sprintf("%v", value)
	}

	return value
}
