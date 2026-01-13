package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TaskEditor interface {
	Load() error

	Add(name string, due string) error                         // Create
	Get() []Task                                               // Read
	Remove(id int16) error                                     // Delete
	RemoveAll() error                                          // Delete all
	Update(id int16, name string, due string, done bool) error // Update v2

	SetFile(name string) error
	Save() error
}

type TaskManager struct {
	db *gorm.DB
} // receiver of TaskEditor

func (tm *TaskManager) Load() error {
	connStr := os.Getenv("SUPABASE_URL")

	if connStr == "" {
		return errors.New("SUPABASE_URL environment variable is not set")
	}

	fmt.Println(connStr)

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	tm.db = db
	sqlDB, _ := tm.db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := db.AutoMigrate(&Task{}); err != nil {
		return fmt.Errorf("automigrate gagal: %v", err)
	}

	return nil
}

func (tm *TaskManager) Add(name string, due string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Task name must not be empty")
	}
	if strings.TrimSpace(due) == "" {
		return errors.New("Task due date must not be empty")
	}

	task := Task{
		Name: name,
		Due:  due,
	}

	res := tm.db.Create(&task)
	return res.Error
}

func (tm *TaskManager) Get() ([]Task, error) {
	var tasks []Task
	if tm.db == nil {
		return nil, fmt.Errorf("db connection is not initialized")
	}

	result := tm.db.Order("id asc").Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (tm *TaskManager) Remove(id int16) error {
	result := tm.db.Delete(&Task{}, id)
	return result.Error
}

func (tm *TaskManager) RemoveAll() error {
	result := tm.db.Where("1 = 1").Delete(&Task{})
	return result.Error
}

func (tm *TaskManager) Update(id int16, name string, due string, done bool) error {
	var task Task
	if err := tm.db.First(&task, id).Error; err != nil {
		return fmt.Errorf("task #%d not found", id)
	}

	task.Name = name
	task.Due = due
	task.Done = done

	result := tm.db.Save(&task)
	return result.Error
}
