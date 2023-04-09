package initialize

import (
	"gopkg.in/yaml.v2"
	"hot-cache/config"
	"hot-cache/global"
	"os"
)

func InitGlobal(path string) error {
	// init config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	global.AppConfig = new(config.Config)
	err = yaml.Unmarshal(yamlFile, global.AppConfig)
	if err != nil {
		return err
	}

	return nil
}
