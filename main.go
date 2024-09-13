package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {

	// root and general handler
	http.HandleFunc("/api/", HandleGenericPaths)

	// server is playing ping-pong
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/api/todos", HandleTodoListAndPost)

	http.HandleFunc("/", HomePage)
	http.Handle("/public/js/", http.StripPrefix("/public/js/", http.FileServer(http.Dir("public/js"))))
	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	// http.Handle("/css", http.FileServer(http.Dir("public/css/")))

	http.ListenAndServe(":3000", nil)
	fmt.Println("Server started...")
}

func HandleTodoListAndPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		HandleTodoList(w, r)
	} else if r.Method == "POST" {
		HandleTodoPost(w, r)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(405)
		w.Write(createResponseBody("Method not supported"))
	}
}

func HandleTodoList(w http.ResponseWriter, r *http.Request) {
	todosBytes, err := json.Marshal(getTodos())
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(todosBytes)
}

func HandleTodoPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var todo Todo
	fmt.Printf("r.Body: %#v\n", body)
	err = json.Unmarshal(body, &todo)
	if err != nil {
		fmt.Printf("err when parsing r.Body: %#v\n", err)
		fmt.Printf("Now trying r.Body parsing via url.ParseQuery, means \"content-type\": \"x-www-form-urlencoded\"\n")
		queryValues, err := url.ParseQuery(string(body))
		if err != nil {
			fmt.Printf("error when parsing r.Body via url.ParseQuery: %#v\n", err)
		} else {
			for key, values := range queryValues {
				if key == "text" {
					todo.Text = values[0]
				}
			}
		}
	}
	if createTodo(todo) {
		w.WriteHeader(201)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(createResponseBody("Todo already exist"))
	}
}

func HandleGetDeleteAndUpdateById(w http.ResponseWriter, r *http.Request) {
	splits := strings.Split(r.URL.Path, "/api/todos/")
	if r.Method == "GET" {
		HandleGetById(w, r, splits[1])
	} else if r.Method == "PUT" {
		HandleUpdateById(w, r)
	} else if r.Method == "DELETE" {
		HandleDeleteById(w, r, splits[1])
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(405)
		w.Write(createResponseBody("Method not supported"))
	}
}

func HandleGetById(w http.ResponseWriter, r *http.Request, todoId string) {
	todo := getTodoById(todoId)
	if todo != nil {
		todoBytes, err := json.Marshal(todo)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(todoBytes)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(createResponseBody("Todo item not found"))
	}
}

func HandleUpdateById(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		panic(err)
	}
	if updateTodo(todo) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(createResponseBody("Todo item updated succeefully"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(createResponseBody("Todo item not found"))
	}
}

func HandleDeleteById(w http.ResponseWriter, r *http.Request, todoId string) {
	if deleteTodoById(todoId) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(createResponseBody("Todo item deleted succeefully"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(createResponseBody("Todo item not found"))
	}
}

func HandleGenericPaths(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/todos/") {
		HandleGetDeleteAndUpdateById(w, r)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write(createResponseBody("Resource not found"))
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("public/index.html"))
	w.Header().Add("Content-Type", "text/html")
	err := templ.Execute(w, getTodos())
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}
}
