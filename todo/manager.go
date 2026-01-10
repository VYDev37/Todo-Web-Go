package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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
	taskList []Task
	fileName string
} // receiver of TaskEditor

func (tm *TaskManager) Load() error {
	data, err := os.ReadFile(tm.fileName)

	if errors.Is(err, os.ErrNotExist) {
		empty := []map[string]Task{}
		data, err := json.Marshal(empty)
		if err != nil {
			return fmt.Errorf("error when trying to write new file: %v", err.Error())
		}

		os.WriteFile(tm.fileName, data, 0644)
		fmt.Printf("File %s does not exist, creating one...\n", tm.fileName)

		return nil
	}
	if err != nil {
		return fmt.Errorf("error when trying to read file: %v", err.Error())
	}

	if err := json.Unmarshal(data, &tm.taskList); err != nil {
		return fmt.Errorf("error when trying to parse data: %v", err.Error())
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

	var newId int16 = 1
	if len(tm.taskList) > 0 {
		newId = tm.taskList[len(tm.taskList)-1].ID + 1
	}

	tm.taskList = append(tm.taskList, Task{
		ID:   newId, // for safety
		Name: name,
		Due:  due,
		Done: false,
	})

	if err := tm.Save(); err != nil {
		return err
	}

	fmt.Println("Saved tasks!")
	return nil
}

func (tm *TaskManager) Get() []Task {
	return tm.taskList
}

func (tm *TaskManager) Remove(id int16) error {
	for i, task := range tm.taskList {
		if task.ID == id {
			tm.taskList = append(tm.taskList[:i], tm.taskList[i+1:]...) // remove element on array based on index
			if err := tm.Save(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("task #%d not found", id)
}

func (tm *TaskManager) RemoveAll() error {
	tm.taskList = tm.taskList[:0]
	// sve
	if err := tm.Save(); err != nil {
		return err
	}

	return nil
}

func (tm *TaskManager) Update(id int16, name string, due string, done bool) error {
	for i, task := range tm.taskList {
		if task.ID == id {
			tm.taskList[i].Name = name
			tm.taskList[i].Due = due
			tm.taskList[i].Done = done

			if err := tm.Save(); err != nil {
				return err
			}

			return nil
		}
	}
	return fmt.Errorf("task #%d not found", id)
}

func (tm *TaskManager) SetFile(name string) error {
	if strings.TrimSpace(name) == "" || !strings.HasSuffix(name, ".json") {
		return errors.New("file extension must be json")
	}

	tm.fileName = name
	return nil
}

func (tm *TaskManager) Save() error {
	data, err := json.MarshalIndent(tm.taskList, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshalling json: %v", err.Error())
	}

	if err := os.WriteFile(tm.fileName, data, 0644); err != nil {
		return fmt.Errorf("error when triyng to write json into file: %v", (err.Error()))
	}

	return nil
}
