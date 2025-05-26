package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"Название"`
	Description string `json:"Описание"`
	Completed   bool   `json:"Завершено"`
}

var nextID int
var tasks = []Task{
	{ID: 1, Title: "Первая задача", Description: "Описание задачи", Completed: false},
	{ID: 2, Title: "Вторая задача", Description: "Ещё одно описание", Completed: true},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Ошибка при чтении задачи", http.StatusBadRequest)
		return
	}
	nextID++
	newTask.ID = nextID
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Ошибка при обновлении", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Title = updated.Title
			tasks[i].Description = updated.Description
			tasks[i].Completed = updated.Completed
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Задача не найдена", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Задача не найдена", http.StatusNotFound)
}

func main() {
	// Проверка наличия папки frontend
	if _, err := os.Stat("../frontend"); os.IsNotExist(err) {
		log.Fatal("❌ Папка ./frontend не найдена")
	}

	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasks(w, r)
		case http.MethodPost:
			createTask(w, r)
		case http.MethodPut:
			updateTask(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("✅ Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
