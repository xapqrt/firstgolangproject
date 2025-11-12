//message types for websocket communication
//trying to figure out json encoding in go

package main

import (
	"encoding/json"
	"fmt"
	"time"
)



type Message struct {
	Type   string `json:"type"`
	TaskID string `json:"taskId"`
	Title  string `json:"title"`
	Status string `json:"status"`
}



func broadcastTask(msg Message) {
	
	var data []byte
	var err error
	
	data, err = json.Marshal(msg)
	
	if err != nil {
		fmt.Println("json marshal error:", err)
		return
	}
	
	fmt.Println("broadcasting:", string(data))
	
	hub.broadcast <- data
}




//handles incoming messages from clients

func handleClientMessage(data []byte) {
	
	var msg Message
	
	err := json.Unmarshal(data, &msg)
	
	if err != nil {
		fmt.Println("json unmarshal error:", err)
			fmt.Println("raw data:", string(data))
		return
	}
	
	fmt.Println("got message type:", msg.Type)
	
	
	if msg.Type == "task_add" {
		
		var newTask Task
		newTask.ID = msg.TaskID
			newTask.Title = msg.Title
		newTask.Status = msg.Status
		newTask.CreatedAt = time.Now()
		
		globalBoard.addTask(newTask)
		
		broadcastTask(msg)
		
		fmt.Println("task added and broadcasted")
		
	} else if msg.Type == "task_move"{
		
		//update task status in board
		var success bool
		success = globalBoard.moveTask(msg.TaskID, msg.Status)
		
		if success {
			broadcastTask(msg)
			fmt.Println("task moved, broadcasted")
		} else {
			fmt.Println("task move failed, task not found")
		}
	}
}
