# GoBoard - Real-Time Kanban Board

yoo this is a real-time multiplayer kanban board i built with go and websockets. its pretty sick ngl

## Features

- **Real-time sync** - multiple ppl can use it at the same time and see changes instantly
- **Drag & drop** - just drag tasks between columns (todo, in progress, done)
- **Due dates** - add due dates to tasks, overdue ones show up in red
- **Delete tasks** - click the X to delete (syncs across all clients)
- **Dark mode** - toggle with the moon icon, saves ur preference
- **Task counters** - shows how many tasks in each column
- **Auto-save** - everything saves to tasks.json automatically

## Tech Stack

- **Backend**: Go 1.25.4 with gorilla/websocket
- **Frontend**: Vanilla HTML/CSS/JS (no frameworks cuh)
- **Storage**: JSON file (tasks.json)
- **WebSocket**: Hub pattern with channels for real-time broadcast

## How to Run

1. make sure u got go installed
2. install dependencies:
   ```bash
   go mod download
   ```
3. run the server:
   ```bash
   go run .
   ```
4. open ur browser to `http://localhost:8080`
5. open multiple tabs to test the real-time sync

## Project Structure

- `main.go` - http server, task struct, board logic
- `hub.go` - websocket hub for managing clients
- `client.go` - individual websocket client handler
- `websocket.go` - websocket upgrade handler
- `message.go` - message handling (add/move/delete tasks)
- `persistence.go` - save/load tasks from json file
- `initial_sync.go` - sends existing tasks to new clients
- `board.html` - kanban ui with drag-drop
- `board.css` - styles (dark mode included)

## Notes

- tasks persist across server restarts (saved in tasks.json)
- dark mode preference saved in browser localStorage
- all websocket messages use json format
- mutex locks prevent race conditions on task operations

## Future Ideas (maybe)

- user accounts/auth
- task comments
- search/filter
- switch to actual database
- deploy it somewhere

---


