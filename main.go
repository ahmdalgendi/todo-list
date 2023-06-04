package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func fetchTasksFromDB() ([]Task, error) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Done) // Add other fields as per your database schema
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func convertTasksToJSON(tasks []Task) ([]byte, error) {
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}

	return tasksJSON, nil
}

func writeJSONResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func main() {
	var err error
	db, err = sql.Open("mysql", "<connection>")
	if err != nil {
		log.Print("error opening database")
		log.Fatal(err)
	}

	r := mux.NewRouter()
	// Set up routes and start the server
	r.HandleFunc("/tasks", tasksHandler)
	r.HandleFunc("/tasks/{id}", oneTaskHandler)

	log.Fatal(http.ListenAndServe(":8080", r))

}

func oneTaskHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	switch request.Method {
	case http.MethodGet:
		getOneTaskHandler(writer, request, id)
	case http.MethodPut:
		updateOneTaskHandler(writer, request, id)
	case http.MethodDelete:
		deleteOneTaskHandler(writer, request, id)
	default:
		http.Error(writer, request.Method, http.StatusNotImplemented)
	}
}

func deleteOneTaskHandler(writer http.ResponseWriter, request *http.Request, id string) {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(writer, []byte("Task deleted successfully"))
}

func updateOneTaskHandler(writer http.ResponseWriter, request *http.Request, id string) {
	var task Task
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE tasks SET  done = ? WHERE id = ?", task.Done, id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err = fetchOneTaskFromDB(id)
	if err != nil {
		http.Error(writer, "Task not found", http.StatusNotFound)
		return
	}
	taskJSON, err := convertTasksToJSON([]Task{task})
	writeJSONResponse(writer, taskJSON)
}

func getOneTaskHandler(writer http.ResponseWriter, request *http.Request, id string) {
	task, err := fetchOneTaskFromDB(id)
	if err != nil {
		http.Error(writer, "Task not found", http.StatusNotFound)
		return
	}

	tasks := []Task{task}
	taskJSON, err := convertTasksToJSON(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(writer, taskJSON)
}

func fetchOneTaskFromDB(id string) (Task, error) {
	var task Task
	err := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Name, &task.Done)
	if err != nil {

		return Task{}, err
	}

	return task, nil
}
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	//http.Error(w, r.Method, http.StatusNotImplemented)
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	case http.MethodPost:
		addTasksHandler(w, r)
	default:
		http.Error(w, r.Method, http.StatusNotImplemented)
	}

}

func getTasksHandler(w http.ResponseWriter, request *http.Request) {
	tasks, err := fetchTasksFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasksJSON, err := convertTasksToJSON(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, tasksJSON)
}

func addTasksHandler(writer http.ResponseWriter, request *http.Request) {
	//read the request body
	var task Task
	err := json.NewDecoder(request.Body).Decode(&task)

	if err != nil {
		http.Error(writer, "here error"+err.Error(), http.StatusInternalServerError)
		return
	}
	//insert into db
	_, err = db.Exec("INSERT INTO tasks(name) VALUES(?)", task.Name)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writeJSONResponse(writer, []byte("success"))
}
