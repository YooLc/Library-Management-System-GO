package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

var DB *gorm.DB

func resetDatabase() {
	if DB == nil {
		logrus.Panic("restting database before connecting to it")
	}
	logrus.Debug("resetting database")
	DB.Migrator().DropTable(&Book{}, &Card{}, &Borrow{})
	DB.AutoMigrate(&Book{}, &Card{}, &Borrow{})
}

func ConnectDatabase(config DatabaseConfig) {
	logrus.Info("connecting to database")
	dsn := fmt.Sprint(config.User, ":", config.Password, "@tcp(", config.Host, ":", config.Port, ")/", config.Database, "?charset=utf8mb4&parseTime=True&loc=Local")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Panic("failed to connect database")
	}
	logrus.Info("connected to database ", DB.Name())
	resetDatabase()
}
