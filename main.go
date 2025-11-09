//alr chat dis is my first go project, tryna make a kanban board, lets see if it works!!




package main


import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

//i think this is how u make a struct?



type Task struct {
	ID        string
	Title     string
	Status    string
	CreatedAt time.Time
}

// gonna store sm tasks in memory for now, might learn to make databses later

type Board struct {
	tasks map[string][]Task
	mu    sync.Mutex
}

var globalBoard *Board
var task_counter int = 0
var hub *Hub

//setup function here ig

func initBoard() {

	fmt.Println("setting up board...")

	// i feel like this is kinda like python wdyt

	globalBoard = &Board{
		tasks: make(map[string][]Task),
	}

	//nned three columns for now

	globalBoard.tasks["todo"] = []Task{}
	globalBoard.tasks["doing"] = []Task{}
	globalBoard.tasks["done"] = []Task{}

	fmt.Println("Board ready!!!")

}


//adding tasks to the board

func (b *Board) addTask(task Task) {

	b.mu.Lock()
	defer b.mu.Unlock()

	var column string
	column = task.Status

	//appening to the right coloumn

	b.tasks[column] = append(b.tasks[column], task)

	fmt.Println("new task added: ", task.Title)
	fmt.Println("current tasks in", column, ":", len(b.tasks[column]))

}






//getting all the tasks
//it returns a map i thinkw wll ckeck

func (b *Board) getAllTasks() map[string][]Task {

	b.mu.Lock()
	defer b.mu.Unlock()

	//making a copy cuz the code camp guy told it

	result := make(map[string][]Task)

	for status, task_list := range b.tasks {

		//copying each list

		var copied_list []Task
		copied_list = append([]Task{}, task_list...)

		result[status] = copied_list
	}

	return result
}



//this is the server fheck func

func healthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("health check hit")

	w.Write([]byte("server is running\n"))
	w.Write([]byte("looks good rn\n"))
}





//main page handlerr

func rootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("someone visited the root page")
	//i hopt it doesnt spam

	//getting all tasks

	all_tasks := globalBoard.getAllTasks()

	//sm html

	w.Write([]byte("<h1>GoBoard</h1>\n"))
	w.Write([]byte("<p>my kanban board project!</p>\n"))

	//shows task counts

	for status, tasks := range all_tasks {

		line := fmt.Sprintf("<p>%s: %d tasks</p>\n", status, len(tasks))
		w.Write([]byte(line))
	}
}







//tryna make a function that creates task id

func makeTaskID() string {

	task_counter = task_counter + 1

	id := fmt.Sprintf("task-%d", task_counter)

	return id
}




func main() {

	fmt.Println("Starting GoBoard server!    ..........")

	//initialiseing the board

	initBoard()

	//setup websocket hub
	hub = newHub()
	go hub.run()

	fmt.Println("websocket hub running in background")

	//adding sm test tasks rn to see if it works

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

	//settin up routes now

	fmt.Println("alr lemme start up routes")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", wsHandler)
 
	fmt.Println("websocket endpoint: /ws")

	//starting the server

	var port string
	port = ":8080"

	fmt.Println("")
	fmt.Println("Server statring on http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	//lets hope this stops server

	err := http.ListenAndServe(port, nil)

	if err != nil {

		fmt.Println("errorr", err)
		fmt.Println("something went wrong, maybe check console")

	}
}