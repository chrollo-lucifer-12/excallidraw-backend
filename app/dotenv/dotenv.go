package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DATABASE_URL string
	PORT         string
}

func NewEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("PORT")
	databaseUrl := os.Getenv("DATABASE_URL")
	env := Env{}
	env.DATABASE_URL = databaseUrl
	env.PORT = port
	return &env, nil
}
