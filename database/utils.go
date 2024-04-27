package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

type APIResult struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

var DB *gorm.DB

func ResetDatabase() {
	if DB == nil {
		logrus.Panic("resting database before connecting to it")
	}
	logrus.Debug("resetting database")
	DB.Migrator().DropTable(&Book{}, &Card{}, &Borrow{})
	DB.AutoMigrate(&Book{}, &Card{}, &Borrow{})
}

func initDatabase() {
	if DB == nil {
		logrus.Panic("initing database before connecting to it")
	}
	logrus.Debug("initing database")
	if DB.Migrator().HasTable(&Book{}) {
		logrus.Debug("table book exists")
	} else {
		logrus.Debug("table book not exists")
		DB.AutoMigrate(&Book{})
	}
	if DB.Migrator().HasTable(&Card{}) {
		logrus.Debug("table card exists")
	} else {
		logrus.Debug("table card not exists")
		DB.AutoMigrate(&Card{})
	}
	if DB.Migrator().HasTable(&Borrow{}) {
		logrus.Debug("table borrow exists")
	} else {
		logrus.Debug("table borrow not exists")
		DB.AutoMigrate(&Borrow{})
	}
}

func ConnectDatabase(config Config) {
	logrus.Info("connecting to database")
	dsn := fmt.Sprint(config.User, ":", config.Password, "@tcp(", config.Host, ":", config.Port, ")/", config.Database, "?charset=utf8mb4&parseTime=True&loc=Local")
	var err error

	logLevel := logger.Default.LogMode(logger.Info)
	switch config.LogLevel {
	case "silent":
		logLevel = logger.Default.LogMode(logger.Silent)
	case "error":
		logLevel = logger.Default.LogMode(logger.Error)
	case "warn":
		logLevel = logger.Default.LogMode(logger.Warn)
	case "info":
	default:
		logLevel = logger.Default.LogMode(logger.Info)
	}
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logLevel,
	})
	if err != nil {
		logrus.WithError(err).Panic("failed to connect database")
	}
	logrus.Info("connected to database ", DB.Name())
	//ResetDatabase()
	initDatabase()
}
