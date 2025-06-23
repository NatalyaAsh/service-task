package models

import "sync"

type TasksBlock struct {
	Tasks  map[int]Task
	LastId int
	Mu     sync.RWMutex
}

type Task struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Created string `json:"created"`
	Lasting int    `json:"lasting"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

type ResponseId struct {
	ID string `json:"id"`
}

const (
	TaskLasting    = 300 // 5 минут
	CreateFormate  = "02.01.2006 15:04:05"
	StatusCreated  = "CREATED"
	StatusFinished = "FINISHED"
)
