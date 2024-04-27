package main

import (
	"fmt"
	"library-management-system/database"
	"library-management-system/server"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var DB *gorm.DB

type AppConfig struct {
	Server   server.Config   `yaml:"server"`
	Database database.Config `yaml:"database"`
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		PadLevelText:  true,
		FullTimestamp: true,
	})

	file, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println("Failed to open config file: ", err)
		return
	}
	defer file.Close()

	var config AppConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		fmt.Println("Failed to parse config file: ", err)
		return
	}

	database.ConnectDatabase(config.Database)
	server.InitServer(config.Server)
}
