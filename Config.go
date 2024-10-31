package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	RedisDataBaseAddress      string
	InterfaceName             string
	ListOfBlockedCountryCodes []string
}

var configInstance *Config

func ReadConfigFromFile() (Config, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	path := configPath + "/goPrometheusPacketMonitoring"
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return Config{}, err
		}

		config := Config{
			RedisDataBaseAddress: "CHANGE_ME",
			InterfaceName:        "CHANGE_ME",
		}

		configAsBytes, err := json.Marshal(config)
		if err != nil {
			return Config{}, err
		}

		os.WriteFile(path+"/config.json", configAsBytes, 0777)
		configInstance = &config
		return config, nil
	}

	configAsBytes, err := os.ReadFile(path + "/config.json")
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(configAsBytes, &config)
	if err != nil {
		return Config{}, err
	}

	configInstance = &config
	return config, nil
}
