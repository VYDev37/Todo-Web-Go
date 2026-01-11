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

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&Task{})
	tm.db = db

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

func (tm *TaskManager) Get() []Task {
	var tasks []Task
	tm.db.Order("id asc").Find(&tasks)
	return tasks
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
