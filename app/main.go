package main

import (
	"fmt"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	fileupload "github.com/chrollo-lucifer-12/excallidraw-backend/app/filleupload"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/server"
)

func main() {
	var err error
	env, err := dotenv.NewEnv()
	if err != nil {
		fmt.Println("error loading env", err)
	}
	dbOpts := db.DBOpts{
		Env: env,
	}
	database, err := db.NewDB(dbOpts)
	if err != nil {
		fmt.Println("error loading db", err)
	}
	uploadServiceOpts := fileupload.UploadServiceOpts{
		Env: env,
	}
	minio := fileupload.NewUploadService(uploadServiceOpts)
	if err != nil {
		fmt.Println("error loading minio", err)
	}
	serverOpts := server.ServerOpts{
		Env:          env,
		Database:     database,
		UploadClient: minio,
	}
	server := server.NewServer(serverOpts)
	server.Start()
}
