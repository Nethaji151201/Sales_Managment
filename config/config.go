package config

import (
	"customersales/apps/models"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	globalConfig *models.Config
	once         sync.Once
)

func LoadGlobalConfig(filename string) {
	once.Do(func() {
		var cfg models.Config
		if _, err := toml.DecodeFile(filename, &cfg); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		globalConfig = &cfg
		log.Println("toml loaded successfully")
	})
}

func GetConfig() *models.Config {
	if globalConfig == nil {
		log.Fatal("no data fount in toml file")
	}
	return globalConfig
}
