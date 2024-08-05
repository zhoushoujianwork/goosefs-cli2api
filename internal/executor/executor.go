package executor

import (
	"goosefs-cli2api/config"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type TaskStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type TaskRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

var (
	tasks      = make(map[string]*exec.Cmd)
	tasksMutex sync.RWMutex
)

func StartTaskManager() {
	// 可以扩展任务清理逻辑等
}

// 只允许内部调用，不允许外部传入所有指令，防止执行影响系统的指令
func addTask(req TaskRequest) (string, error) {
	taskID := uuid.New().String()
	cmd := exec.Command(req.Command, req.Args...)
	tasksMutex.Lock()
	tasks[taskID] = cmd
	tasksMutex.Unlock()

	outputPath := strings.TrimSuffix(*config.Config.OutputDir, "/") + "/" + taskID + ".txt"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	cmd.Stdout = outFile
	cmd.Stderr = outFile

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	go func() {
		cmd.Wait()
		outFile.Close()
	}()

	return taskID, nil
}

func GetTaskStatus(taskID string) (TaskStatus, error) {
	tasksMutex.RLock()
	defer tasksMutex.RUnlock()
	cmd, exists := tasks[taskID]
	if !exists {
		return TaskStatus{}, os.ErrNotExist
	}
	return TaskStatus{ID: taskID, Status: cmd.ProcessState.String()}, nil
}

func GetTaskOutput(taskID string) (string, error) {
	outputPath := strings.TrimSuffix(*config.Config.OutputDir, "/") + "/" + taskID + ".txt"
	content, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
