/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/api"
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"

	"github.com/alibabacloud-go/tea/tea"
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
		config.Init(debug)
		if debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		go restartTask()
		r := gin.Default()
		api.RegisterRoutes(r) // 注册API路由
		log.Infof("api server start on http://localhost:%d", port)
		if err := r.Run(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
			panic(err)
		}
	},
}

// 防止任务异常中断有没有成功的任务，继续跑起来
func restartTask() {
	tasks, err := config.DB.GetGoosefsTask(models.FilterGoosefsTaskRequest{})
	if err != nil {
		panic(err)
	}
	reqs := make(map[string]models.GooseFSRequest)
	// 重新组装原始 req 请求
	for _, task := range tasks {
		if task.TaskName == nil ||
			(task.Action != models.GFSForceLoad && task.Action != models.GFSDistributeLoad) {
			log.Debugf("task_name and taskids are empty or action is not GFSForceLoad or GFSDistributeLoad, skip checkTasksIsFinished")
			continue
		}
		if executor.GetCmdStatus(tea.StringValue(task.ExitCode)) == models.TaskStatusRunning {
			if _, ok := reqs[*task.TaskName]; !ok {
				reqs[*task.TaskName] = models.GooseFSRequest{
					TaskName: task.TaskName,
					Action:   task.Action,
					Path:     []*string{task.Path},
				}
			} else {
				req := reqs[*task.TaskName]
				req.Path = append(req.Path, task.Path)
				reqs[*task.TaskName] = req
			}
		}
	}

	for _, req := range reqs {
		// 过滤出已经成功执行的任务
		tasks, err := config.DB.GetGoosefsTask(models.FilterGoosefsTaskRequest{
			TaskName: req.TaskName,
			Action:   &req.Action,
		})
		if err != nil {
			log.Errorf("get tasks error: %s", err)
			continue
		}
		newPath := make([]*string, 0)
		for _, task := range tasks {
			if executor.GetCmdStatus(tea.StringValue(task.ExitCode)) == models.TaskStatusSuccess {
				for _, path := range req.Path {
					if *path == *task.Path {
						log.Debugf("task has success, skip task: %s", tea.Prettify(task))
					} else {
						newPath = append(newPath, path)
					}
				}
			}
		}

		if len(newPath) == 0 {
			log.Debugf("all tasks has success, skip task: %s", tea.Prettify(req))
			continue
		}
		req.Path = newPath
		// 重新执行
		switch req.Action {
		case models.GFSDistributeLoad:
			resp, err := executor.DistrubuteLoad(req)
			if err != nil {
				log.Errorf("restart distrubuteLoad error: %s", err)
				continue
			} else {
				log.Infof("restart distrubuteLoad success: %s", tea.Prettify(resp))
			}

		case models.GFSForceLoad:
			err := executor.ForceLoad(req)
			if err != nil {
				log.Errorf("restart forceLoad error: %s", err)
				continue
			} else {
				log.Infof("restart forceLoad success: %s", tea.Prettify(req))
			}
		default:
			log.Errorf("restart unknown action: %s", tea.Prettify(req))
		}

	}

	log.Infof("reload task& restart finished and nu:%d", len(reqs))

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
