package main

import (
	"goosefs-cli2api/internal/executor"
	"testing"
)

// TEST DistrubuteLoad
func TestDistrubuteLoad(t *testing.T) {
	taskID, err := executor.DistrubuteLoad("/data-datalake/deltalake/aaa.db/bbb/")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(taskID)
}
