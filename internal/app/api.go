package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"service-task/internal/config"
	"service-task/internal/models"
	"service-task/internal/tasks"
	"strconv"
	"time"
	// "service-demo/internal/config"
	// "service-demo/internal/database/pgsql"
	// "service-demo/internal/database/redis"
	// modeldb "service-demo/internal/models"
)

// type Meta struct {
// 	Total   int `json:"total"`
// 	Removed int `json:"removed"`
// 	Limit   int `json:"limit"`
// 	Offset  int `json:"offset"`
// }

// type StructGetGoods struct {
// 	Meta  Meta             `json:"meta"`
// 	Goods *[]modeldb.Goods `json:"goods"`
// }

func Init(mux *http.ServeMux, cfg *config.Config) {
	mux.HandleFunc("POST /task", PostTask)
	mux.HandleFunc("GET /task", GetTask)
	mux.HandleFunc("DELETE /task", DeleteTask)
}

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	msg, _ := json.Marshal(data)
	io.Writer.Write(w, msg)
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	//slog.Info("Api.PostTask")
	var buf bytes.Buffer
	var task models.Task

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, models.ResponseErr{Error: "ошибка конвертации данных"}, http.StatusBadRequest)
		return
	}

	fmt.Printf("task: %v\n", task)
	id := tasks.PostTask(task)
	// Запускаем горутину на время выполнения задачи
	go WorkTask(w, id)

	writeJson(w, models.ResponseId{ID: strconv.Itoa(id)}, http.StatusOK)

}

func GetTask(w http.ResponseWriter, r *http.Request) {
	//slog.Info("Api.GetTask")

	IdRaw := r.URL.Query().Get("id")
	if IdRaw == "" {
		writeJson(w, models.ResponseErr{Error: "не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(IdRaw)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: "не корректный идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	//slog.Info("Api.GetTask", "Id", id)
	task, err := tasks.GetTask(id)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, task, http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	//slog.Info("Api.DeleteTask")

	IdRaw := r.URL.Query().Get("id")
	if IdRaw == "" {
		writeJson(w, models.ResponseErr{Error: "не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(IdRaw)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: "не корректный идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	//slog.Info("Api.DeleteTask", "Id", id)
	task, err := tasks.GetTask(id)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Status != models.StatusFinished {
		writeJson(w, models.ResponseErr{Error: "задача ещё выполняется, её нельзя удалить"}, http.StatusOK)
		return
	}

	err = tasks.DeleteTask(id)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, "", http.StatusOK)
}

func WorkTask(w http.ResponseWriter, id int) {
	task, err := tasks.GetTask(id)
	if err != nil {
		writeJson(w, models.ResponseErr{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	time.Sleep(time.Duration(task.Lasting) * time.Second)
	tasks.TaskFinished(task.ID)
	slog.Info("Api.WorkTask Finished", "id", task.ID)
}
