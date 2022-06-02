package cmd

import (
	"fmt"
	"github.com/ivankoTut/nltools/helpers"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Сохранение настроек",
	Long:  `Сохранение настроек`,
	Run: func(cmd *cobra.Command, args []string) {
		updateDevServerDir(cmd)
		showConfig(cmd)
		urlShow(cmd)
		urlShowAll(cmd)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().String("dir-dev-server", "", "Путь к папке dev-server")
	configCmd.Flags().BoolP("show", "s", false, "Показать конфиг")
	configCmd.Flags().BoolP("url-show", "u", false, "Показать локальные адреса(с совпадениями по сервисам)")
	configCmd.Flags().BoolP("url-show-all", "S", false, "Показать все локальные адреса")
}

func urlShowAll(cmd *cobra.Command) {
	isShowUrl, _ := cmd.Flags().GetBool("url-show-all")

	if isShowUrl == false {
		return
	}

	var myEnv map[string]string
	myEnv, err := godotenv.Read(viper.GetString("envPath"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Переменная", "Ссылка"})

	for variable, value := range myEnv {
		table.Append([]string{variable, value})
	}

	table.Render()
}

func urlShow(cmd *cobra.Command) {
	isShowUrl, _ := cmd.Flags().GetBool("url-show")

	if isShowUrl == false {
		return
	}

	cfg := &helpers.Config{}
	if err := helpers.ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		panic(err)
	}
	var myEnv map[string]string
	myEnv, err := godotenv.Read(viper.GetString("envPath"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Сервис", "Ссылка"})

	for key, _ := range cfg.Application {
		service := strings.Replace(key, "_", ".", -1)
		for _, value := range myEnv {
			if strings.Contains(value, service) {
				table.Append([]string{service, value})
			}
		}
	}

	table.Render()
}

func showConfig(cmd *cobra.Command) {
	isShow, _ := cmd.Flags().GetBool("show")

	if isShow == false {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Наименование", "Значение"})
	data := map[string]string{}
	homeDir, _ := os.UserHomeDir()

	helpers.ReadYamlToConfig(homeDir+"/.nltools", data)
	for key, value := range data {
		table.Append([]string{key, value})
	}

	table.Render()
}

func updateDevServerDir(cmd *cobra.Command) {
	dir, _ := cmd.Flags().GetString("dir-dev-server")
	if dir == "" {
		return
	}
	configPath := dir + "/config.yaml"
	envPath := dir + "/.env"

	if helpers.FolderExist(configPath) != true {
		fmt.Println(fmt.Sprintf("%s не найден", configPath))
		os.Exit(0)
	}

	if helpers.FolderExist(envPath) != true {
		fmt.Println(fmt.Sprintf("%s не найден", envPath))
		os.Exit(0)
	}

	viper.Set("rootDit", dir)
	viper.Set("configPath", configPath)
	viper.Set("envPath", envPath)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Путь добавлен в конфиг")
}
