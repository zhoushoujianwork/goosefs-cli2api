package db

import (
	"goosefs-cli2api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

func NewSqliteDB(dbFile string, debug bool) DB {
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if debug {
		dbConfig.Logger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(sqlite.Open(dbFile), dbConfig)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.GoosefsTask{})
	if err != nil {
		panic(err)
	}

	return DB{db: db}
}

func (d *DB) CreateGoosefsTask(taskID string, req models.GoosefsTaskRequest) error {
	return d.db.Create(req.ToGoosefsTask(taskID)).Error
}

func (d *DB) UpdateGoosefsTask(taskID string, req models.UpdateGoosefsTaskRequest) error {
	return d.db.Model(&models.GoosefsTask{}).Where("id = ?", taskID).Updates(req).Error
}

func (d *DB) GetGoosefsTask(filter models.FilterGoosefsTaskRequest) ([]models.GoosefsTask, error) {
	sql := d.db.Model(&models.GoosefsTask{})
	if filter.TaskID != nil {
		sql = sql.Where("id = ?", *filter.TaskID)
	}
	if filter.TaskName != nil {
		sql = sql.Where("task_name = ?", *filter.TaskName)
	}
	if filter.Action != nil {
		sql = sql.Where("action = ?", *filter.Action)
	}
	var tasks []models.GoosefsTask
	err := sql.Find(&tasks).Error
	return tasks, err
}

// func (d *DB) CreateGoosefsPathStatus(req models.CreatePathStatusRequest) error {
// 	return d.db.Create(req.ToPathStatus()).Error
// }

// func (d *DB) UpdateGoosefsPathStatus(statusID string, req models.UpdatePathStatusRequest) error {
// 	return d.db.Model(&models.GoosefsPathStatus{}).Where("id = ?", statusID).Updates(req).Error
// }

// func (d *DB) GetGoosefsPathStatus(filter models.FilterPathStatusRequest) ([]models.GoosefsPathStatus, error) {
// 	sql := d.db.Model(&models.GoosefsPathStatus{})
// 	if filter.TaskID != nil {
// 		sql = sql.Where("task_id = ?", *filter.TaskID)
// 	}
// 	if filter.Path != nil {
// 		sql = sql.Where("path = ?", *filter.Path)
// 	}
// 	if filter.TaskName != nil {
// 		sql = sql.Where("task_name = ?", *filter.TaskName)
// 	}
// 	var statuses []models.GoosefsPathStatus
// 	err := sql.Find(&statuses).Error
// 	return statuses, err
// }
