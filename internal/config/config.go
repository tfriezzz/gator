// Package config reads/writes json
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type State struct {
	Config Config
}

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// fmt.Printf("home: %s\n", homeDir)
	gatorConfig := fmt.Sprintf("%s/%s", homeDir, configFileName)
	return gatorConfig, nil
}

func write(cfg Config) error {
	v, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	e := os.WriteFile(filePath, v, 0o644)
	if e != nil {
		return err
	}
	return nil
}

func (c Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(c)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Read() Config {
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err = json.Unmarshal(content, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
