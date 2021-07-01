package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Tabler interface {
	TableName() string
}

var Db *gorm.DB
var path = "todo.db"

func Init() error {
	newLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	var err error
	Db, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: newLogger,
	})

	return err
}
