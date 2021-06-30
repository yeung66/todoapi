package main

import (
	"fmt"
	database "github.com/yeung66/todoapi/internal/db"
	"github.com/yeung66/todoapi/internal/todos"
)

func main() {
	database.Init()
	todoItems := todos.GetAllTodoItems()
	for _, t := range todoItems {
		fmt.Println(t.Title)
	}
}
