package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/asymptoter/practice-backend/base/ctx"
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
	SmtpHost          string `yaml:"smtpHost"`
	Port              int    `yaml:"port"`
	Account           string `yaml:"account"`
	Password          string `yaml:"password"`
	ActivationMessage string `yaml:"activationMessage"`
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
	context := ctx.Background()
	file, err := ioutil.ReadFile(pwd[:len(pwd)-1] + "/../../config/config.yml")
	if err != nil {
		panic("config.Init failed at ioutil.ReadFile " + err.Error())
	}

	if err := yaml.Unmarshal(file, Configurations); err != nil {
		panic("config.Init failed at yaml.Unmarshal " + err.Error())
	}

	temp := Configurations[*ENV]
	b1, _ := json.Marshal(temp)
	b2, _ := json.Marshal(Value)
	if string(b1) == string(b2) {
		return
	}
	Value = temp
	context.Info("Configuration updated.")
}
