package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	tasks map[string][]*Context
	mu    sync.Mutex
}

type Context struct {
	ch       chan string
	IsFinish bool
}

func NewService() *Service {
	tasks := make(map[string][]*Context)
	return &Service{tasks: tasks}
}
func (svc *Service) addQueue(appId string, ctx *Context) {
	svc.mu.Lock()
	if svc.tasks[appId] == nil {
		svc.tasks[appId] = make([]*Context, 0)
	}

	svc.tasks[appId] = append(svc.tasks[appId], ctx)
	svc.mu.Unlock()
}

func (svc *Service) removeQueue(appId string, ctx *Context) {
	svc.mu.Lock()
	newTasks := make([]*Context, 0)
	if _, ok := svc.tasks[appId]; ok {
		for _, ele := range svc.tasks[appId] {
			if ele == ctx {
				continue
			} else {
				newTasks = append(newTasks, ele)
			}
		}
		svc.tasks[appId] = newTasks
	}
	svc.mu.Unlock()
}

func (svc *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	appId := r.URL.Query().Get("app_id")
	if appId == "" {
		w.WriteHeader(400)
		return
	}
	ctx := &Context{ch, false}
	svc.addQueue(appId, ctx)

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("time out, return")
		svc.removeQueue(appId, ctx)
		w.WriteHeader(304)

		return
	case content := <-ch:
		fmt.Println("changed")
		fmt.Fprintf(w, "%s", content)
		return
	}
}

func (svc *Service) Publish(w http.ResponseWriter, r *http.Request) {
	appId := r.URL.Query().Get("app_id")
	content := r.URL.Query().Get("content")
	svc.mu.Lock()
	defer svc.mu.Unlock()
	tasks := svc.tasks[appId]
	fmt.Printf("tasks len = %d\n ", len(tasks))
	for _, ctx := range tasks {

		ctx.ch <- content
	}
	svc.tasks[appId] = nil
	w.WriteHeader(200)
}

func main() {
	svc := NewService()
	http.HandleFunc("/get_config", svc.GetConfig)
	http.HandleFunc("/publish", svc.Publish)
	http.ListenAndServe("localhost:8081", nil)
}
