package configs

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"
)

func InitConfigs() error {

	viper.SetConfigType("yaml")

	appEnv := getAppEnv()

	conf, err := ioutil.ReadFile(appEnv)
	if err != nil {
		return errors.New("tonic_error.configs.missing_configs_file")
	}

	err = viper.ReadConfig(bytes.NewBuffer(conf))
	if err != nil {
		return errors.New("configs_error.configs.invalid_format")
	}

	return nil

}

func getAppEnv() string {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "./configs/development.yaml"
	}
	return appEnv
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}
