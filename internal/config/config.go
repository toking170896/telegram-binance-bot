package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Config ...
type Config struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	ChatID   int64 `json:"chatID"`
	Db       string `json:"db"`
}

//Init / New / Get
func Init() (*Config, error) {
	conf := &Config{}
	_, err := os.Stat("config.json")

	if os.IsNotExist(err) {
		return nil, err
	}

	b, err := ioutil.ReadFile("config.json")
	if err := json.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
