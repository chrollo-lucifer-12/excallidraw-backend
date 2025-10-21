package main

import (
	"fmt"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/server"
)

func main() {
	var err error
	env, err := dotenv.NewEnv()
	if err != nil {
		fmt.Println("error loading env")
	}
	dbOpts := db.DBOpts{
		Env: env,
	}
	database, err := db.NewDB(dbOpts)
	if err != nil {
		fmt.Println("error loading db")
	}
	serverOpts := server.ServerOpts{
		Env:      env,
		Database: database,
	}
	server := server.NewServer(serverOpts)
	server.Start()
}
