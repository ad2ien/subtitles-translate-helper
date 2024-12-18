package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Test  string `yaml:"test"`
	Input struct {
		Path string `yaml:"path"`
		Lang string `yaml:"lang"`
	} `yaml:"input"`
	Output struct {
		Path string `yaml:"path"`
		Lang string `yaml:"lang"`
	} `yaml:"output"`
	IgnoreSubStartingWithChar  string `yaml:"ignoreSubStartingWithChar"`
	LibreTranslateServicePort  string `yaml:"libreTranslateServicePort"`
	LibreTranslateImageVersion string `yaml:"libreTranslateImageVersion"`
}

var config Config

func GetConfig() *Config {
	return &config
}

func InitConfig(configPath string) *Config {
	buf, error := os.ReadFile(configPath)
	if error != nil {
		logger.Println("ERROR can't read config file")
		logger.Panic(error)
	}

	error = yaml.Unmarshal(buf, &config)
	if error != nil {
		logger.Println("ERROR can't parse config file")
		logger.Panic(error)
	}

	return &config

}

func GetLangArgument() string {
	return config.Input.Lang + "," + config.Output.Lang
}
