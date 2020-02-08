package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	ENV            = flag.String("env", "local", "environment")
	Configurations = map[string]Configuration{}
	Value          = Configuration{}
)

type ServerConfiguration struct {
	Address string             `yaml:"address"`
	Email   EmailConfiguration `yaml:"email"`
}

type EmailConfiguration struct {
	Address          string `yaml:"address"`
	Port             int    `yaml:"port"`
	OfficialAccount  string `yaml:"officialAccount"`
	OfficialPassword string `yaml:"officialPassword"`
}

type MysqlConfiguration struct {
	Address         string `yaml:"address"`
	DatabaseName    string `yaml:"databaseName"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	ConnectionRetry int    `yaml:"connectionRetry"`
}

type RedisConfiguration struct {
	Address string `yaml:"address"`
}

type Configuration struct {
	Server ServerConfiguration `yaml:"server"`
	MySQL  MysqlConfiguration  `yaml:"mysql"`
	Redis  RedisConfiguration  `yaml:"redis"`
}

func Init() {
	file, err := ioutil.ReadFile("../../config/config.yml")
	if err != nil {
		panic("ioutil.ReadFile failed " + err.Error())
	}

	if err := yaml.Unmarshal(file, Configurations); err != nil {
		panic("yaml.Unmarshal failed " + err.Error())
	}

	Value = Configurations[*ENV]
}
