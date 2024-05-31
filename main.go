package main

import (
	"log/slog"
	"sync"
	"time"
)

type Task struct {
	ID          int
	Description string
	Action      func()
}

// Worker function to process tasks from the tasks channel and send completion signals
func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		task.Action()
		slog.Info("Worker finished task", slog.Int("worker_id", id), slog.Int("task_id", task.ID), slog.String("description", task.Description))
	}
}

// TaskScheduler struct to manage task scheduling
type TaskScheduler struct {
	maxConcurrentTasks int
	tasks              chan Task
	wg                 sync.WaitGroup
}

// NewTaskScheduler creates a new TaskScheduler
func NewTaskScheduler(maxConcurrentTasks int) *TaskScheduler {
	return &TaskScheduler{
		maxConcurrentTasks: maxConcurrentTasks,
		tasks:              make(chan Task),
	}
}

// Start initializes the worker pool
func (ts *TaskScheduler) Start() {
	for w := 1; w <= ts.maxConcurrentTasks; w++ {
		ts.wg.Add(1)
		go worker(w, ts.tasks, &ts.wg)
	}
}

// AddTask adds a new task to the scheduler
func (ts *TaskScheduler) AddTask(task Task) {
	ts.tasks <- task
}

// Wait waits for all tasks to complete and then closes the tasks channel
func (ts *TaskScheduler) Wait() {
	close(ts.tasks)
	ts.wg.Wait()
}

func main() {
	// Create a task scheduler with a maximum of 3 concurrent tasks
	scheduler := NewTaskScheduler(3)
	scheduler.Start()

	tasks := []Task{
		{ID: 1, Description: "Task 1", Action: func() { time.Sleep(time.Second * 1) }},
		{ID: 2, Description: "Task 2", Action: func() { time.Sleep(time.Second * 2) }},
		{ID: 3, Description: "Task 3", Action: func() { time.Sleep(time.Second * 3) }},
		{ID: 4, Description: "Task 4", Action: func() { time.Sleep(time.Second * 4) }},
		{ID: 5, Description: "Task 5", Action: func() { time.Sleep(time.Second * 5) }},
	}

	for _, task := range tasks {
		scheduler.AddTask(task)
	}

	scheduler.Wait()
	slog.Info("All tasks completed")
}
