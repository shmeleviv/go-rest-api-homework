package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/go-chi/chi/v5"
)

type NetIP struct {
	//CimClass              map[string]string `json:"CimClass"`
	//CimInstanceProperties []string          `json:"CimInstanceProperties"`
	//CimSystemProperties   map[string]string `json:"CimSystemProperties"`
	InterfaceAlias string `json:"InterfaceAlias"`
	IPv4Address    string `json:"IPv4Address"`
}
type Operation struct {
	Op   string   `json:"op"`
	Path []string `json:"path"`
}

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getIPAddr(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	netip := []NetIP{}
	var op Operation
	var c *exec.Cmd
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	fmt.Println(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &op); err != nil {
		fmt.Println(op)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(op)
	c = exec.Command("powershell", "/C", op.Path[0])
	ipconfig, _ := c.CombinedOutput()
	if err = json.Unmarshal(ipconfig, &netip); err != nil {
		fmt.Println(netip)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("NetIP Struct:", netip)
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	resp, _ := json.Marshal(netip)
	w.Write(resp)
}

// Ниже напишите обработчики для каждого эндпоинта
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)
	r.Get("/ipAddr", getIPAddr)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
