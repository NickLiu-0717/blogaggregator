package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func getConfigPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homedir + "/" + configFilename, nil
}

func Read() (Config, error) {
	filepath, err := getConfigPath()
	if err != nil {
		fmt.Println(err)
	}
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	fullpath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
