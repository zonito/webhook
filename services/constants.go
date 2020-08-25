package services

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	PushoverKey  string
	TrelloKey    string
	TrelloSecret string
	TeleToken    string
}

func getConfig() Config {
	content, _ := ioutil.ReadFile("services/keys.json")
	var conf Config
	json.Unmarshal(content, &conf)
	return conf
}

var config = getConfig()

var pushoverKey = config.PushoverKey
var trelloKey = config.TrelloKey
var trelloSecret = config.TrelloSecret
var teleToken = config.TeleToken
