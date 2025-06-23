package tasks

import (
	"fmt"
	"log/slog"
	"service-task/internal/models"
	"time"
)

var Tasks models.TasksBlock

func Init() {
	Tasks.Tasks = make(map[int]models.Task)
	Tasks.LastId = 0
}

func PostTask(task models.Task) int {
	if task.Lasting == 0 {
		task.Lasting = models.TaskLasting
	}
	Tasks.Mu.Lock()
	Tasks.LastId++
	task.ID = Tasks.LastId
	task.Created = time.Now().UTC().Format(models.CreateFormate)
	task.Status = models.StatusCreated

	Tasks.Tasks[task.ID] = task
	Tasks.Mu.Unlock()
	slog.Info("Tasks.PostTask", "task", task)
	return Tasks.LastId
}

func GetTask(id int) (*models.Task, error) {
	slog.Info("Tasks.GetTask", "id", id)
	Tasks.Mu.RLock()
	task, exist := Tasks.Tasks[id]
	if exist {
		Tasks.Mu.RUnlock()
		return &task, nil
	}
	Tasks.Mu.RUnlock()
	return nil, fmt.Errorf("задача не найдена")
}

func DeleteTask(id int) error {
	slog.Info("Tasks.DeleteTask", "id", id)
	Tasks.Mu.Lock()
	_, exist := Tasks.Tasks[id]
	if exist {
		delete(Tasks.Tasks, id)
		Tasks.Mu.Unlock()
		return nil
	}
	Tasks.Mu.Unlock()
	return fmt.Errorf("такой задачи не существует")
}

func TaskFinished(id int) {
	Tasks.Mu.Lock()
	task, exist := Tasks.Tasks[id]
	if exist {
		Tasks.Tasks[id] = models.Task{ID: task.ID, Status: models.StatusFinished, Created: task.Created, Lasting: task.Lasting}
	}
	Tasks.Mu.Unlock()
}
