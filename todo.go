package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
)

func writeTodos(todos []Todo) {
	todosBytes, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("todos.json", todosBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func getTodos() (todos []Todo) {
	todoBytes, err := os.ReadFile("todos.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(todoBytes, &todos)
	if err != nil {
		panic(err)
	}
	return todos
}

func getTodoById(id string) *Todo {
	var todo Todo
	todos := getTodos()

	for i, _ := range todos {
		if todos[i].Id == id {
			todo = todos[i]
			break
		}
	}
	if todo.Id == "" {
		return nil
	}
	return &todo
}

func deleteTodoById(id string) bool {
	var found bool = false
	var tempTodos []Todo
	for _, todo := range getTodos() {
		if todo.Id != id {
			tempTodos = append(tempTodos, todo)
		} else {
			found = true
		}
	}
	if found {
		writeTodos(tempTodos)
	}
	return found
}

func updateTodo(todo Todo) bool {
	var found bool = false
	todos := getTodos()

	for i, _ := range todos {
		if todos[i].Id == todo.Id {
			todos[i].Text = todo.Text
			found = true
			break
		}
	}

	if found {
		writeTodos(todos)
	}

	return found
}

func createTodo(todo Todo) bool {
	if todo.Text == "" {
		return false
	}
	todo.Id = uuid.NewString()
	todos := getTodos()
	for _, t := range todos {
		if strings.Compare(t.Text, todo.Text) == 0 {
			return false
		}
	}
	todos = append(todos, todo)
	writeTodos(todos)
	return true
}
