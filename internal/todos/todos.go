package todos

import (
	"github.com/yeung66/todoapi/internal/db"
)

type TodoItem struct {
	Id          int
	Title       string
	Content     string
	Checked     bool
	CreatedTime string `gorm:"column:createdTime"`
	UpdatedTime string `gorm:"column:updatedTime"`
}

func (TodoItem) TableName() string {
	return "TodoItems"
}

func GetAllTodoItems() []TodoItem {
	var todos []TodoItem
	db := database.Db
	db.Find(&todos)
	return todos
}

func GetTodoItem(id int) TodoItem {
	var todo TodoItem
	database.Db.First(&todo, id)
	return todo
}

func (t TodoItem) Save() error {
	result := database.Db.Create(&t)
	return result.Error
}

func (t *TodoItem) Update(attrs map[string]interface{}) error {
	result := database.Db.Model(t).Updates(attrs)
	return result.Error
}

func DeleteTodoItem(id int) TodoItem {
	var todo TodoItem
	database.Db.First(&todo, id)
	database.Db.Delete(&todo)
	return todo
}

func GetTodoItemsByTimeRange(start, end *string) []TodoItem {
	if start == nil && end == nil {
		return GetAllTodoItems()
	}

	var todos []TodoItem
	if end == nil {
		database.Db.Find(&todos, "createdTime < ?", *start)
	} else if start == nil {
		database.Db.Find(&todos, "createdTime > ?", *end)
	} else {
		database.Db.Find(&todos, "createdTime > ? AND createdTime < ?", *end, *start)
	}

	return todos
}
