package helpers

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func FolderExist(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	return false
}

func ReadYamlToConfig(path string, cnf interface{}) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(yamlFile, cnf)
	}
	return err
}

func SaveYamlConfig(config *Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(viper.GetString("configPath"), data, 0777); err != nil {
		return err
	}

	return nil
}
