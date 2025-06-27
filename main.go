package main

import (
	"fmt"
	"go_final_project/go_final_project/pkg/db"
	"log"
	"os"

	"go1f/pkg/server"
)
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env файл не найден")
	}
	port := os.Getenv("TODO_PORT")

	DB, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", tools.AuthMiddleware(handlers.TaskActionHandler))
	http.HandleFunc("/api/tasks", tools.AuthMiddleware(handlers.GetAllTasks))
	http.HandleFunc("/api/task/done", tools.AuthMiddleware(handlers.TaskDone))
	http.HandleFunc("/api/signin", handlers.SignInHandler)

	fmt.Println("Сервер запущен на порте:", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
func taskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		tasks_service.AddTaskHandler(w, r, db)
	case http.MethodGet:
		tasks_service.TasksHandler(w, r, db)
	case http.MethodPut:
		tasks_service.EditTaskHandler(w, r, db)
	case http.MethodDelete:
		tasks_service.DeleteTaskHandler(w, r, db)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
