package cmd

import (
	"fmt"
	"github.com/ivankoTut/go-alerts"
	"github.com/ivankoTut/nltools/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	printTable = true
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Основные команды по управлению сервисами",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printTable, _ = cmd.Flags().GetBool("table")

		showAllServices(cmd)
		showDevServices(cmd)
		showStableServices(cmd)
		restartsContainers(cmd)
		disableService(cmd)
		enableService(cmd)
		stopService(cmd)
		startService(cmd)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().BoolP("show-dev", "d", false, "Показать сервисы работающие локально")
	appCmd.Flags().BoolP("show-all", "a", false, "Показать все сервисы")
	appCmd.Flags().BoolP("show-stable", "s", false, "Показать сервисы не работающие локально")
	appCmd.Flags().BoolP("table", "t", true, "Не печатать в табличном виде")
	appCmd.Flags().BoolP("restart", "r", false, "Перезапустить контейнеры (./nlrun stop && ./nlrun)")
	appCmd.Flags().BoolP("stop", "S", false, "Остановить сервисы (./nlrun stop)")
	appCmd.Flags().BoolP("start", "T", false, "Запустить сервисы (./nlrun)")
	appCmd.Flags().String("disable", "", "Переключить сервис в режим stable")
	appCmd.Flags().String("enable", "", "Переключить сервис в режим dev")
}

func startService(cmd *cobra.Command) {
	isStop, _ := cmd.Flags().GetBool("start")
	if isStop == false {
		return
	}
	alerts.CreateBlock("Запуск сервисов ....", "", helpers.NoticeTheme)
	helpers.StartServices()
	alerts.CreateBlock("Сервисы успешно запущены", "", helpers.NoticeTheme)
}

func stopService(cmd *cobra.Command) {
	isStop, _ := cmd.Flags().GetBool("stop")
	if isStop == false {
		return
	}
	alerts.CreateBlock("Остановка сервисов ....", "", helpers.NoticeTheme)
	helpers.StopServices()
	alerts.CreateBlock("Сервисы успешно остановлены ....", "", helpers.NoticeTheme)
}

func enableService(cmd *cobra.Command) {
	service, _ := cmd.Flags().GetString("enable")
	if service == "" {
		return
	}

	if !helpers.IsExistService(service, helpers.ModeStable) {
		fmt.Println(fmt.Sprintf("Сервис %s не найден", service))
		os.Exit(1)
	}

	if err := helpers.EnableService(service); err != nil {
		fmt.Println(err)
	}

	if !helpers.IsExistService(service, helpers.ModeDev) {
		fmt.Println(fmt.Sprintf("Не удалось изменить статус сервиса %s на %s", service, helpers.ModeDev))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Статус сервиса %s измене на %s", service, helpers.ModeDev))
	fmt.Println("Рестартим сервисы")
	helpers.RestartServices()
}

func disableService(cmd *cobra.Command) {
	service, _ := cmd.Flags().GetString("disable")
	if service == "" {
		return
	}

	if !helpers.IsExistService(service, helpers.ModeDev) {
		fmt.Println(fmt.Sprintf("Сервис %s не найден", service))
		os.Exit(1)
	}

	if err := helpers.DisableService(service); err != nil {
		fmt.Println(err)
	}

	if !helpers.IsExistService(service, helpers.ModeStable) {
		fmt.Println(fmt.Sprintf("Не удалось изменить статус сервиса %s на %s", service, helpers.ModeStable))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Статус сервиса %s измене на %s", service, helpers.ModeStable))
	fmt.Println("Рестартим сервисы")
	helpers.RestartServices()
}

func restartsContainers(cmd *cobra.Command) {
	isRestart, _ := cmd.Flags().GetBool("restart")
	if isRestart == false {
		return
	}

	helpers.RestartServices()
}

func showAllServices(cmd *cobra.Command) {
	isShow, _ := cmd.Flags().GetBool("show-all")
	if isShow == false {
		return
	}

	cfg := &helpers.Config{}
	if err := helpers.ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		log.Panicln(err)
	}

	var table *tablewriter.Table
	if printTable {
		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Сервис", "Режим"})
	}

	for key, value := range cfg.Application {
		if value.Mode == helpers.ModeDev {
			if printTable {
				table.Rich([]string{key, value.Mode}, []tablewriter.Colors{{tablewriter.FgWhiteColor, tablewriter.Normal, tablewriter.BgGreenColor}, {tablewriter.FgWhiteColor, tablewriter.Normal, tablewriter.BgGreenColor}})
			} else {
				fmt.Println(fmt.Sprintf("%s: %s", key, value.Mode))
			}
		} else {
			if printTable {
				table.Append([]string{key, value.Mode})
			} else {
				fmt.Println(fmt.Sprintf("%s: %s", key, value.Mode))
			}
		}
	}

	if printTable {
		table.Render()
	}
}

func showStableServices(cmd *cobra.Command) {
	isShow, _ := cmd.Flags().GetBool("show-stable")
	if isShow == false {
		return
	}

	cfg := &helpers.Config{}
	if err := helpers.ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		log.Panicln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Сервис", "Режим"})
	for key, value := range cfg.Application {
		if value.Mode == helpers.ModeDev {
			continue
		}
		table.Append([]string{key, value.Mode})
	}

	table.Render()
}

func showDevServices(cmd *cobra.Command) {
	isShow, _ := cmd.Flags().GetBool("show-dev")
	if isShow == false {
		return
	}

	cfg := &helpers.Config{}
	if err := helpers.ReadYamlToConfig(viper.GetString("configPath"), cfg); err != nil {
		log.Panicln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Сервис", "Режим"})
	for key, value := range cfg.Application {
		if value.Mode != helpers.ModeDev {
			continue
		}
		table.Append([]string{key, value.Mode})
	}

	table.Render()
}
