package main

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

// TEST DistrubuteLoad
func TestDistrubuteLoad(t *testing.T) {
	taskID, err := executor.DistrubuteLoad(models.GooseFSRequest{
		TaskName: tea.String("test"),
		Path: []*string{
			tea.String("/data-datalake/deltalake/aaa.db/bbb/"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(taskID)
}
