package main

import (
	"context"
	"embed"
	"encoding/json"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

const tasksFile = "tasks.json"

type Task struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowShow(ctx)
	runtime.LogInfo(a.ctx, "App has started")
}

func (a *App) FetchTasks() []Task {
	file, err := os.Open(tasksFile)
	if err != nil {
		return []Task{}
	}
	defer file.Close()
	var tasks []Task
	decoder := json.NewDecoder(file)
	decoder.Decode(&tasks)
	return tasks
}

func (a *App) SaveTasks(tasks []Task) {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	_ = os.WriteFile(tasksFile, data, 0644)
}

func (a *App) AddTask(text string) []Task {
	if text == "" {
		return a.FetchTasks()
	}
	tasks := a.FetchTasks()
	tasks = append(tasks, Task{Text: text, Completed: false})
	a.SaveTasks(tasks)
	return tasks
}

func (a *App) RemoveTask(index int) []Task {
	tasks := a.FetchTasks()
	if index < 0 || index >= len(tasks) {
		return tasks
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	a.SaveTasks(tasks)
	return tasks
}

func (a *App) ToggleTask(index int) []Task {
	tasks := a.FetchTasks()
	if index < 0 || index >= len(tasks) {
		return tasks
	}
	tasks[index].Completed = !tasks[index].Completed
	a.SaveTasks(tasks)
	return tasks
}

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:       "TodoApp",
		Width:       800,
		Height:      600,
		Frameless:   false,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
