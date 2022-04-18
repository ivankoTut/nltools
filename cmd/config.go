/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/ivankoTut/nltools/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Сохранение настроек",
	Long:  `Сохранение настроек`,
	Run: func(cmd *cobra.Command, args []string) {
		updateDevServerDir(cmd)
		showConfig(cmd)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().String("dir-dev-server", "", "Путь к папке dev-server")
	configCmd.Flags().BoolP("show", "s", false, "Показать конфиг")
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
