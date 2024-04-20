package server

import (
	"fmt"
	"library-management-system/database"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   Config                  `yaml:"server"`
	Database database.DatabaseConfig `yaml:"database"`
}

func TestMain(m *testing.M) {
	file, err := os.Open("../config.yaml")
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
	m.Run()
}

func TestBookRegister(t *testing.T) {
	b0 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", Publish_year: 2023,
		Author: "Mike", Price: 188.88, Stock: 10,
	}
	assert.Equal(t, StoreBook(b0).Ok, true)

	/* Not allowed to create duplicated records */
	b1 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", Publish_year: 2023,
		Author: "Mike", Price: 188.88, Stock: 5,
	}
	b2 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", Publish_year: 2023,
		Author: "Mike", Price: 99.99, Stock: 10,
	}
	assert.Equal(t, StoreBook(b1).Ok, false)
	assert.Equal(t, StoreBook(b2).Ok, false)
}
