package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type WcaiConfiguration struct {
	ClientId     string
	ClientSecret string
	DbHost       string
	DbName       string
	DbUser       string
	DbPassword   string
}

func GetConfiguration() WcaiConfiguration {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("error:", err)
		panic("Unable to open config file")
	}
	decoder := json.NewDecoder(file)
	configuration := WcaiConfiguration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		panic("Unable to get configuration")
	}

	return configuration
}
