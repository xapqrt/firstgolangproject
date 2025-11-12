//send initial tasks to new clients when they connect
//so they dont start with empty board

package main

import (
	"encoding/json"
	"fmt"
)


//sends all existing tasks to a newly connected client
func sendInitialTasks(c *Client) {
	
	fmt.Println("sending initial tasks to new client...")
	
	var allTasks map[string][]Task
	allTasks = globalBoard.getAllTasks()
	
	//loop through each status column
	for status, tasks := range allTasks {
		
		fmt.Printf("sending %d tasks from %s column\n", len(tasks), status)
		
		for _, task := range tasks {
			
			//create message for each task
			var msg Message
			msg.Type = "task_add"
			msg.TaskID = task.ID
			msg.Title = task.Title
				msg.Status = task.Status
			
			//marshal to json
			var data []byte
			var err error
			data, err = json.Marshal(msg)
			
			if err != nil {
				fmt.Println("json marshal errorr:", err)
				continue
			}
			
			//send directly to this client only
			//not broadcasting to everyone
			select {
			case c.send <- data:
				fmt.Println("sent task:", task.ID, "to new client")
			default:
				fmt.Println("client send buffer full, skipping task")
			}
		}
	}
	
	fmt.Println("finished sending initial tasks")
}
