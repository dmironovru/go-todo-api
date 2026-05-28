package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "github.com/jackc/pgx/v5"
)

type Task struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

var db *pgx.Conn

func main() {
    var err error
    connString := "postgres://todo_user:todo123@localhost:5432/todo_db"
    db, err = pgx.Connect(context.Background(), connString)
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }
    defer db.Close(context.Background())
    log.Println("Connected to PostgreSQL")

    http.HandleFunc("GET /api/tasks", corsMiddleware(getTasks))
    http.HandleFunc("POST /api/tasks", corsMiddleware(createTask))
    http.HandleFunc("PUT /api/tasks/{id}", corsMiddleware(updateTask))
    http.HandleFunc("DELETE /api/tasks/{id}", corsMiddleware(deleteTask))
    http.HandleFunc("DELETE /api/tasks", corsMiddleware(clearTasks))

    log.Println("Сервер запущен на :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query(context.Background(), "SELECT id, title, completed FROM tasks ORDER BY id")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    tasks := []Task{}
    for rows.Next() {
        var t Task
        rows.Scan(&t.ID, &t.Title, &t.Completed)
        tasks = append(tasks, t)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title string `json:"title"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
        http.Error(w, "Invalid title", http.StatusBadRequest)
        return
    }

    var task Task
    err := db.QueryRow(context.Background(),
        "INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, completed",
        req.Title).Scan(&task.ID, &task.Title, &task.Completed)
    
    if err != nil {
        log.Println("Database error:", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid id", http.StatusBadRequest)
        return
    }

    var req struct {
        Completed bool `json:"completed"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    _, err = db.Exec(context.Background(),
        "UPDATE tasks SET completed=$1 WHERE id=$2", req.Completed, id)
    if err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid id", http.StatusBadRequest)
        return
    }

    _, err = db.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1", id)
    if err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func clearTasks(w http.ResponseWriter, r *http.Request) {
    _, err := db.Exec(context.Background(), "DELETE FROM tasks")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
