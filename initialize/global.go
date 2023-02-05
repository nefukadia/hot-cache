package initialize

import (
	"gopkg.in/yaml.v2"
	"hot-cache/cache"
	"hot-cache/config"
	"hot-cache/global"
	"os"
)

const ConfigFile = "./config.yaml"

func InitGlobal() error {
	// init config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}
	global.AppConfig = new(config.Config)
	err = yaml.Unmarshal(yamlFile, global.AppConfig)
	if err != nil {
		return err
	}

	// init cache
	global.AppCache = cache.New()

	return nil
}
