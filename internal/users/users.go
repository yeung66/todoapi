package users

import (
	"github.com/dgryski/trifles/uuid"
	database "github.com/yeung66/todoapi/internal/db"
)

type User struct {
	Id          int
	CreatedTime string `gorm:"column:createdTime"`
	Username    string
	Password    string
	Token       string
}

func (User) TableName() string {
	return "Users"
}

func GetUserByToken(token string) User {
	var user User
	database.Db.Where("token = ?", token).First(&user)
	return user
}

func GetUserByLogin(username, password string) User {
	var user User
	results := database.Db.Where(&User{Password: password, Username: username}).First(&user)
	if results.RowsAffected == 1 {
		user.Token = uuid.UUIDv4()
	}
	return user
}
