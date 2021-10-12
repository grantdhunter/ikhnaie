package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Port                       string
	Mapbox_api_key             string
	Github_oauth_url           string
	Github_oauth_clientid      string
	Github_oauth_client_secret string
}

func (c Configuration) String() string {
	str, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", str)

}

var Config Configuration

func init() {
	if _, err := toml.DecodeFile("./config.toml", &Config); err != nil {
		log.Fatal(err)
		return
	}

}
