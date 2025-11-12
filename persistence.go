//save and load tasks from file so they dont disappear on restart
//using json encoding to write to tasks.json file

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const tasksFile = "tasks.json"

//saves all tasks to file
func saveTasks() error {
	
	fmt.Println("saving tasks to file...")
	
	var allTasks map[string][]Task
	allTasks = globalBoard.getAllTasks()
	
	//marshal to json with indentation so its readable
	var data []byte
	var err error
	data, err = json.MarshalIndent(allTasks, "", "  ")
	
	if err != nil {
		fmt.Println("json marshal errorr:", err)
		return err
	}
	
	//write to file
	err = os.WriteFile(tasksFile, data, 0644)
	
	if err != nil {
		fmt.Println("file write error:", err)
		return err
	}
	
	fmt.Println("tasks saved sucessfully")
	return nil
}

//loads tasks from file on startup
func loadTasks() error {
	
	fmt.Println("loading tasks from file...")
	
	//check if file exists
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		fmt.Println("no tasks file found, starting fresh")
		return nil
	}
	
	//read file
	var data []byte
	var err error
	data, err = os.ReadFile(tasksFile)
	
	if err != nil {
		fmt.Println("file read error:", err)
		return err
	}
	
	//unmarshal json
	var loadedTasks map[string][]Task
	err = json.Unmarshal(data, &loadedTasks)
	
	if err != nil {
		fmt.Println("json unmarshal errorr:", err)
		return err
	}
	
	//add tasks to board
	globalBoard.mu.Lock()
	globalBoard.tasks = loadedTasks
	globalBoard.mu.Unlock()
	
	//update task counter to highest task id
	for _, taskList := range loadedTasks {
		for _, task := range taskList {
			//extract number from task id
			//task id is like "task-5" so we parse the number part
			fmt.Println("loaded task:", task.ID, task.Title)
		}
	}
	
	fmt.Println("tasks loaded sucessfully")
	return nil
}
