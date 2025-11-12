//alr chat dis is my first go project, tryna make a kanban board, lets see if it works!!


package main


import (
	"fmt"
	"net/http"
	"sync"
	"time"
)


type Task struct {
	ID        string
	Title     string
	Status    string

	CreatedAt time.Time
	DueDate   string    `json:"dueDate,omitempty"`
}



type Board struct {
	tasks map[string][]Task
	mu    sync.Mutex
}



var globalBoard *Board
var task_counter int = 0
var hub *Hub

func initBoard() {




	fmt.Println("setting up board...")

	globalBoard = &Board{
		tasks: make(map[string][]Task),
	}


	globalBoard.tasks["todo"] = []Task{}

	globalBoard.tasks["doing"] = []Task{}
	globalBoard.tasks["done"] = []Task{}

	fmt.Println("Board ready!!!")

}


func (b *Board) addTask(task Task) {

	b.mu.Lock()
	defer b.mu.Unlock()

	var column string
	column = task.Status

	b.tasks[column] = append(b.tasks[column], task)

	fmt.Println("new task added: ", task.Title)
	fmt.Println("current tasks in", column, ":", len(b.tasks[column]))

}



func (b *Board) getAllTasks() map[string][]Task {

	b.mu.Lock()

	defer b.mu.Unlock()

	result := make(map[string][]Task)

	for status, task_list := range b.tasks {

		var copied_list []Task



		copied_list = append([]Task{}, task_list...)

		result[status] = copied_list
	}

	return result
}







func (b *Board) moveTask(taskID string, newStatus string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	//find 


	for oldStatus, taskList := range b.tasks {
		for i, task := range taskList {
			if task.ID == taskID {



				//remove 
				// 
				
				b.tasks[oldStatus] = append(taskList[:i], taskList[i+1:]...)
				
				//update 



				task.Status = newStatus
				b.tasks[newStatus] = append(b.tasks[newStatus], task)
				
				fmt.Println("moved task", taskID, "from", oldStatus, "to", newStatus)
				return true
			}
		}
	}
	
	fmt.Println("task not found:", taskID)
	return false
}

//delete task from board
func (b *Board) deleteTask(taskID string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	for status, taskList := range b.tasks {
		for i, task := range taskList {
			if task.ID == taskID {
				b.tasks[status] = append(taskList[:i], taskList[i+1:]...)
				fmt.Println("deleted task", taskID, "from", status)
				return true
			}
		}
	}
	
	fmt.Println("task not found for deletion:", taskID)
	return false
}


func healthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("health check hit")

	w.Write([]byte("server is running\n"))



	w.Write([]byte("looks good rn\n"))
}



func rootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("someone visited the root page")

	all_tasks := globalBoard.getAllTasks()

	w.Write([]byte("<h1>GoBoard</h1>\n"))
	w.Write([]byte("<p>my kanban board project!</p>\n"))

	for status, tasks := range all_tasks {



		line := fmt.Sprintf("<p>%s: %d tasks</p>\n", status, len(tasks))
		w.Write([]byte(line))



	}
}


func boardHandler(w http.ResponseWriter, r *http.Request){

	fmt.Println("serving board.html")



	
	// w.Write([]byte("board.html") )
	
	



	if r.URL.Path == "/board.css" {
		http.ServeFile(w, r, "board.css")
		return
	}
	
	http.ServeFile(w, r, "board.html")
}





func makeTaskID() string {

	task_counter = task_counter + 1

	id := fmt.Sprintf("task-%d", task_counter)

	return id
}




func main() {

	fmt.Println("Starting GoBoard server!    ..........")

	initBoard()
	
	//load tasks from file if it exists
	err := loadTasks()
	if err != nil {
		fmt.Println("couldnt load tasks, starting fresh")
	}

	
	hub = newHub()
	go hub.run()

	fmt.Println("websocket hub running in background")

	fmt.Println("adding test tasks....")

	task1 := Task{
		ID:        makeTaskID(),
		Title:     "learn golang",
		Status:    "doing",
		CreatedAt: time.Now(),
	}

	globalBoard.addTask(task1)

	task2 := Task{
		ID:        makeTaskID(),
		Title:     "build websocket tmr, learning",
		Status:    "todo",
		CreatedAt: time.Now(),
	}

	globalBoard.addTask(task2)

	var task3 Task
	task3.ID = makeTaskID()
	task3.Title = "make it look cool"
	task3.Status = "todo"


	task3.CreatedAt = time.Now()



	globalBoard.addTask(task3)

	fmt.Println("test tasks added!!")





	fmt.Println("alr lemme start up routes")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/board", boardHandler)
	http.HandleFunc("/board.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "board.css")
	})



	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", wsHandler)
 
	fmt.Println("websocket endpoint: /ws")
	fmt.Println("board ui: /board")




	var port string
	port = ":8080"

	fmt.Println("")
	fmt.Println("Server statring on http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	var serverErr error
	serverErr = http.ListenAndServe(port, nil)

	if serverErr != nil {



		fmt.Println("errorr", serverErr)
		fmt.Println("something went wrong, maybe check console")

	}
}