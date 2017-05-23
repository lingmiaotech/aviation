package aviation

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"
)

type ConfigsClass struct{}

var Configs ConfigsClass

func InitConfigs() error {

	viper.SetConfigType("yaml")

	appEnv := getAppEnv()

	conf, err := ioutil.ReadFile(appEnv)
	if err != nil {
		return errors.New("aviation_error.configs.missing_configs_file")
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

func (c ConfigsClass) Get(key string) interface{} {
	return viper.Get(key)
}

func (c ConfigsClass) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c ConfigsClass) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (c ConfigsClass) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c ConfigsClass) GetString(key string) string {
	return viper.GetString(key)
}

func (c ConfigsClass) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (c ConfigsClass) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func (c ConfigsClass) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (c ConfigsClass) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (c ConfigsClass) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func (c ConfigsClass) IsSet(key string) bool {
	return viper.IsSet(key)
}
