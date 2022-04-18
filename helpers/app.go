package helpers

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os/exec"
)

var (
	ModeStable = "stable"
	ModeDev    = "dev"
)

//Config Представление файла конфигурации в dev-server
type Config struct {
	Gitlab         GitlabConfig                 `yaml:"gitlab"`
	Infrastructure []string                     `yaml:"infrastructure"`
	Application    map[string]ApplicationConfig `yaml:"application"`
}

type ApplicationConfig struct {
	Mode      string `yaml:"mode"`
	RemoteEnv string `yaml:"remoteEnv"`
}

type GitlabConfig struct {
	BaseUrl string `yaml:"baseUrl"`
	Token   string `yaml:"token"`
}

func DisableService(serviceName string) error {
	cfg := &Config{}
	if err := ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		return err
	}
	for key, value := range cfg.Application {
		if key == serviceName {
			current := &value
			current.Mode = ModeStable
			cfg.Application[key] = *current
			if err := SaveYamlConfig(cfg); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func EnableService(serviceName string) error {
	cfg := &Config{}
	if err := ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		return err
	}
	for key, value := range cfg.Application {
		if key == serviceName {
			current := &value
			current.Mode = ModeDev
			cfg.Application[key] = *current
			if err := SaveYamlConfig(cfg); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func IsExistService(serviceName string, mode string) bool {
	cfg := &Config{}
	if err := ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		log.Panicln(err)
	}
	for key, value := range cfg.Application {
		if value.Mode == mode && key == serviceName {
			return true
		}
	}

	return false
}

func RestartServices() {
	fmt.Println("Остановка сервисов ....")
	qwe := exec.Command("nlrun", "stop")
	qwe.Dir = viper.GetString("rootDit")
	_, err := qwe.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Print the output
	fmt.Println("Запуск сервисов ....")

	qwe = exec.Command("nlrun")
	qwe.Dir = viper.GetString("rootDit")
	_, err = qwe.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Сервисы перезапущены")
}
