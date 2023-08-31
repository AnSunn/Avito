package main

import (
	"flag"
	"github.com/AnSunn/ServerUserSegmentation/internal/app/api"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/api.toml", "path to config file (.toml or .env file)")
}

func main() {
	//Flag parsing and add value to configPath var
	flag.Parse()
	//config instance
	config := api.NewConfig()
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Println("!cannot find configs file. Using default values:", err)
	}
	config.BindAddr = os.Getenv("Bind_addr")
	config.LogLevel = os.Getenv("Logger_level")
	config.Store.DatabaseURL = os.Getenv("database_url")

	//server instance
	s := api.New(config)

	//server start
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
