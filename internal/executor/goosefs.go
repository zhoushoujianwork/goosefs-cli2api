package executor

import (
	"fmt"
	"goosefs-cli2api/config"
	"log"
	"time"
)

/*
./bin/goosefs fs distributedLoad --replication 1 /data-datalake/deltalake/aaa.db/bbb/
*/
func DistrubuteLoad(path string) (string, error) {
	return AddTask(TaskRequest{
		Command: *config.Config.Bin,
		Args:    []string{"fs", "distributedLoad", "--replication", "1", path},
	})
}

/*
./bin/goosefs fs loadMetadata -R /data-datalake/deltalake/aaa.db/bbb/
*/
func LoadMetadata(path string) (string, error) {
	return AddTask(TaskRequest{
		Command: *config.Config.Bin,
		Args:    []string{"fs", "loadMetadata", "-R", path},
	})
}

/*
./bin/goosefs fsadmin report
*/
func Report() (string, error) {
	taskid, err := AddTask(TaskRequest{
		Command: *config.Config.Bin,
		Args:    []string{"fsadmin", "report"},
	})
	if err != nil {
		return "", err
	}
	// log.Println("taskid:", taskid)
	// wait for task done
	count := 0
	for {
		count++
		if count > 30 {
			return "", fmt.Errorf("wait for task done timeout")
		}
		// log.Println("get task status:", taskid)
		status, err := GetTaskStatus(taskid)
		if err != nil {
			return "", err
		}
		log.Println("task status:", status.Status)
		if status.Status == "exit status 0" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 读取输出文件
	output, err := GetTaskOutput(taskid)
	if err != nil {
		return "", fmt.Errorf("get task output error: %v", err)
	}
	return output, nil
}
