package users

import (
	database "github.com/yeung66/todoapi/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          int    `gorm:"primaryKey"`
	CreatedTime string `gorm:"column:createdTime"`
	Username    string
	Password    string
	//Token       string
}

type WrongUsernameOrPasswordError struct{}

type PermissionDeniedError struct {
}

type HasNoTodoItemError struct {
}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

func (*PermissionDeniedError) Error() string {
	return "no permission on such operation"
}

func (*HasNoTodoItemError) Error() string {
	return "user has no such todo item"
}

func (User) TableName() string {
	return "Users"
}

func (u *User) Save() error {
	return database.Db.Create(u).Error
}

func GetUserByToken(token string) User {
	var user User
	database.Db.Where("token = ?", token).First(&user)
	return user
}

func GetUserByLogin(username, password string) (User, error) {
	var user User
	results := database.Db.Where(&User{Username: username}).First(&user)
	if results.RowsAffected != 1 || !CheckPasswordHash(password, user.Password) {
		return user, &WrongUsernameOrPasswordError{}
	}
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	var user User
	res := database.Db.Where(&User{
		Username: username,
	}).First(&user)

	return user, res.Error

}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
