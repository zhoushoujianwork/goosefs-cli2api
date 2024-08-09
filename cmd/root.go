/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goosefs-cli2api/internal/api"
	"goosefs-cli2api/internal/executor"
	"os"

	"github.com/xops-infra/noop/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goosefs-cli2api",
	Short: "run goosefs cli to api",
	Long: `by run this tools on goosefs server, you can send request to api server to run goosefs cli.
1. this tools must be run on server, not on client.
2. must has goosefs cli.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		r := gin.Default()
		api.RegisterRoutes(r)          // 注册API路由
		go executor.StartTaskManager() // 启动任务管理器，负责任务的调度和状态管理
		log.Infof("api server start on http://localhost:%d", port)
		if err := r.Run(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	debug bool
	port  int
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goosefs-cli2api.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "api server port")
	if debug {
		os.Setenv("DEBUG", "true")
	}
}
