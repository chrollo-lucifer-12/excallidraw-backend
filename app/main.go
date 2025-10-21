package main

import (
	"fmt"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
)

func main() {
	env, err := dotenv.NewEnv()
	if err != nil {
		fmt.Println("error loading env")
	}
	fmt.Println(env)
}
