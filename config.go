package xmicro

import (
	"fmt"
	"path/filepath"

	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/mofancloud/xmicro/utils"
)

var (
	// AppConfig is the instance of Config, store the config information from file
	AppConfig config.Config

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is json
	appConfigProvider = "json"
)

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	if !utils.FileExists(absConfigPath) {
		return fmt.Errorf("the target config file: %s don't exist", configPath)
	}

	appConfigPath = absConfigPath
	appConfigProvider = adapterName

	return parseConfig(appConfigPath)
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	source := file.NewSource(config.SourceName(appConfigPath))
	AppConfig = config.NewConfig(config.WithSource(source))

	return nil
	/*
		AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
		if err != nil {
			return err
		}
		return assignConfig(AppConfig)
	*/
}
