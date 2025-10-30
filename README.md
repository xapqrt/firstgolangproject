# GoBoard

> a living digital corkboard held together by websocket signals

## what is this

multiplayer kanban board built with raw go
no frameworks, no react, just websockets and chaos

## current status

✅ basic http server running
✅ in-memory task storage (3 columns: todo, doing, done)
✅ mutex-protected board operations
✅ dummy test data

## how to run

```bash
go run main.go
```

then open: `http://localhost:8080`

## the vibe

- chaotic yet minimal
- like a terminal app that became social
- raw Go energy, no sugar
- comments that tell stories

## what's next

- [ ] websocket connections
- [ ] html board view
- [ ] real-time task sync
- [ ] drag and drop
- [ ] multiple clients

## tech stack

- **backend**: Go (net/http)
- **frontend**: HTML + raw JS
- **data**: in-memory map (for now)
- **sync**: websockets (coming soon)

---

built during late night coding sessions
probably has bugs but it works
