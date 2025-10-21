package main

import (
	"fmt"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
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
	_, err = db.NewDB(dbOpts)
	if err != nil {
		fmt.Println("error loading db")
	}

}
