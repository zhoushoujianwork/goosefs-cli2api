package utils

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/models"
	"os"
	"strings"

	"github.com/xops-infra/noop/log"

	"github.com/alibabacloud-go/tea/tea"
)

// 依据传入的 task_name 产生唯一的 task_id，格式类型为 task_name_uuid.txt
func GenerateTaskID(taskName, uid string) string {
	// 清除空格，避免文件名中出现空格
	taskName = strings.ReplaceAll(taskName, " ", "")
	filePath := strings.TrimSuffix(*config.Config.OutputDir, "/") + "/" + taskName + "_" + uid + ".txt"
	log.Debugf("filePath: %s", filePath)
	return filePath
}

// FindFiles 根据taskID或taskName在给定目录下找到所有匹配文件路径
func FindFiles(req models.QueryTaskRequest) ([]string, error) {
	log.Debugf(tea.Prettify(req))

	// 直接依据规则返回确定文件
	if req.TaskID != nil && *req.TaskID != "" && req.TaskName != nil && *req.TaskName != "" {
		return []string{fmt.Sprintf("%s_%s.txt", *req.TaskName, *req.TaskID)}, nil
	}

	var result []string
	// 读取目录内容
	entries, err := os.ReadDir(*config.Config.OutputDir)
	if err != nil {
		log.Panicf("Failed to read directory: %v", err)
	}

	// 遍历目录内的每一项
	for _, entry := range entries {
		// 跳过子目录
		if entry.IsDir() {
			continue
		}
		// log.Infof("FindFiles: ", entry.Name())
		key := entry.Name()
		if (req.TaskID != nil && *req.TaskID != "" && strings.HasSuffix(key, *req.TaskID+".txt")) ||
			(req.TaskName != nil && *req.TaskName != "" && (strings.HasPrefix(key, *req.TaskName))) {
			log.Debugf("FindFiles:%s", key)
			if req.TaskName != nil && *req.TaskName != "" {
				// 排除前缀一致的任务情况。比如 taskA 会匹配到 taskABC这种情况，通过移除 taskname 判断后续是否满足 uuid来判断
				log.Debugf(strings.TrimSuffix(strings.TrimPrefix(key, *req.TaskName+"_"), ".txt"))
				if len(strings.TrimSuffix(strings.TrimPrefix(key, *req.TaskName+"_"), ".txt")) != 36 {
					continue
				}
			}
			result = append(result, key)
		}
	}

	return result, nil
}

// 通过文件名解析出TaskID
func ParseTaskID(fileName string) (string, error) {
	if strings.Count(fileName, "_") == 0 {
		return "", fmt.Errorf("ask admin check might be has '_', invalid file name: %s", fileName)
	}

	parts := strings.Split(fileName, "_")
	return strings.TrimSuffix(parts[len(parts)-1], ".txt"), nil
}
