/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/api"
	"goosefs-cli2api/internal/executor"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xops-infra/noop/log"
)

var port int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "api server",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
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

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "api server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
