package helpers

import (
	"github.com/ivankoTut/go-alerts"
	"github.com/spf13/viper"
	"log"
	"os/exec"
)

var (
	ModeStable  = "stable"
	ModeDev     = "dev"
	NoticeTheme *alerts.Color
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

func init() {
	NoticeTheme, _ = alerts.CreateColor("yellow", "default", []string{"bold"})
	NoticeTheme.
		PrintPaddingBottom(false).
		PrintPaddingTop(false)
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
	alerts.CreateBlock("Остановка сервисов ....", "", NoticeTheme)
	StopServices()

	alerts.CreateBlock("Запуск сервисов ....", "", NoticeTheme)
	StartServices()

	alerts.CreateBlock("Сервисы успешно запущены", "", NoticeTheme)
}

func StartServices() {
	nlRunStart := exec.Command("nlrun")
	nlRunStart.Dir = viper.GetString("rootDit")
	_, err := nlRunStart.Output()

	if err != nil {
		alerts.Error(err.Error())
		return
	}
}

func StopServices() {
	nlRunStop := exec.Command("nlrun", "stop")
	nlRunStop.Dir = viper.GetString("rootDit")
	_, err := nlRunStop.Output()

	if err != nil {
		alerts.Error(err.Error())
		return
	}
}
