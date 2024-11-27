package model

import (
	"gopkg.in/yaml.v3"
	"os"
)

var (
	pwdPath, _ = os.Getwd()
	config     *Config
	// todo 待调整
	path = "/Users/zhangjay/go/solar-metrics/config/solar-metric.yml"
)

type Config struct {
	Solar struct {
		Logrus struct {
			Format       string `yaml:"format"`
			Level        string `yaml:"level"`
			RecordMethod bool   `yaml:"recordMethod"`
		} `yaml:"logrus"`
		Rabbitmq struct {
			Dsl      string `yaml:"dsl"`
			Consumer struct {
				ConnectionNums int `yaml:"connection-nums"`
				ChannelNums    int `yaml:"channel-nums"`
				Queues         []string
				Qos            int    `yaml:"qos"`
				Name           string `yaml:"name"`
			} `yaml:"consumer"`
		} `yaml:"rabbitmq"`
	} `yaml:"solar"`
}

func init() {
	config = &Config{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, config)
}

func GetConfig() *Config {
	return config
}
