package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DATABASE_URL string
}

func NewEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	env := Env{}
	env.DATABASE_URL = databaseUrl

	return &env, nil
}
