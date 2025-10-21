package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DATABASE_URL     string
	PORT             string
	MINIO_ENDPOINT   string
	MINIO_ACCESS_KEY string
	MINIO_SECRET_KEY string
}

func NewEnv() (*Env, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	databaseUrl := os.Getenv("DATABASE_URL")
	minio_endpoint := os.Getenv("MINIO_ENDPOINT")
	minio_access_key := os.Getenv("MINIO_ACCESS_KEY")
	minio_secret_key := os.Getenv("MINIO_SECRET_KEY")
	env := Env{}
	env.DATABASE_URL = databaseUrl
	env.PORT = port
	env.MINIO_ENDPOINT = minio_endpoint
	env.MINIO_ACCESS_KEY = minio_access_key
	env.MINIO_SECRET_KEY = minio_secret_key
	return &env, nil
}
