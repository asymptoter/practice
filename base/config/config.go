package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	ENV            = flag.String("env", "local", "environment")
	WORKDIR        = flag.String("workdir", ".", "working directory")
	Configurations = map[string]Configuration{}
	Value          = Configuration{}
)

type ServerConfiguration struct {
	Address string                     `yaml:"address"`
	Email   OfficialEmailConfiguration `yaml:"email"`
	Testing TestingConfiguration       `yaml:"testing"`
}

type OfficialEmailConfiguration struct {
	SmtpHost string `yaml:"smtpHost"`
	Port     int    `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type TestingConfiguration struct {
	Email EmailConfiguration `yaml:"email"`
}

type EmailConfiguration struct {
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type DatabaseConfiguration struct {
	Address         string `yaml:"address"`
	Port            int    `yaml:"port"`
	DatabaseName    string `yaml:"databaseName"`
	UserName        string `yaml:"userName"`
	Password        string `yaml:"password"`
	ConnectionRetry int    `yaml:"connectionRetry"`
}

type RedisConfiguration struct {
	Address string `yaml:"address"`
}

type Configuration struct {
	Server   ServerConfiguration   `yaml:"server"`
	Database DatabaseConfiguration `yaml:"database"`
	Redis    RedisConfiguration    `yaml:"redis"`
}

func Init(pwd string) {
	file, err := ioutil.ReadFile(pwd[:len(pwd)-1] + "/../../config/config.yml")
	if err != nil {
		panic("ioutil.ReadFile failed " + err.Error())
	}

	if err := yaml.Unmarshal(file, Configurations); err != nil {
		panic("yaml.Unmarshal failed " + err.Error())
	}

	Value = Configurations[*ENV]
}
