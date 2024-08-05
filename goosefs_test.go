package main

import (
	"goosefs-cli2api/internal/executor"
	"testing"
)

// TEST AddTask
func TestAddTask(t *testing.T) {
	req := executor.TaskRequest{
		Command: "ls",
		Args:    []string{"-l"},
	}
	taskID, err := executor.AddTask(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(taskID)
}

// TEST DistrubuteLoad
func TestDistrubuteLoad(t *testing.T) {
	taskID, err := executor.DistrubuteLoad("/data-datalake/deltalake/aaa.db/bbb/")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(taskID)
}
