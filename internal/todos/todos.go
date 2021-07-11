package todos

import (
	"github.com/yeung66/todoapi/internal/db"
	"github.com/yeung66/todoapi/internal/users"
	"gorm.io/gorm"
)

type TodoItem struct {
	gorm.Model
	Id          int `gorm:"primaryKey"`
	Title       string
	Content     string
	Checked     bool
	CreatedTime string `gorm:"column:createdTime"`
	UpdatedTime string `gorm:"column:updatedTime"`
	UserId      int
	User        users.User `gorm:"foreignKey:UserId;references:Id"`
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

func GetUserAllTodoItems(userid int) []TodoItem {
	var todos []TodoItem
	database.Db.Find(&todos, TodoItem{UserId: userid})
	return todos
}

func GetTodoItem(id int) TodoItem {
	var todo TodoItem
	database.Db.First(&todo, id)
	return todo
}

func (t *TodoItem) Save() error {
	result := database.Db.Create(t)
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
		database.Db.Find(&todos, "createdTime >= ?", *start)
	} else if start == nil {
		database.Db.Find(&todos, "createdTime <= ?", *end)
	} else {
		database.Db.Find(&todos, "createdTime <= ? AND createdTime >= ?", *end, *start)
	}

	return todos
}

func GetUserTodoItemsByTimeRange(start, end *string, userid int) []TodoItem {
	if start == nil && end == nil {
		return GetUserAllTodoItems(userid)
	}

	db := database.Db.Where("user_id = ?", userid)
	var todos []TodoItem
	if end == nil {
		db.Find(&todos, "createdTime >= ?", *start)
	} else if start == nil {
		db.Find(&todos, "createdTime <= ?", *end)
	} else {
		db.Find(&todos, "createdTime <= ? AND createdTime >= ?", *end, *start)
	}

	return todos
}

func UserHasTodoItem(userid, id int) bool {
	var todo TodoItem
	database.Db.First(&todo, id)
	return todo.UserId == userid
}
