package fileupload

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadServiceOpts struct {
	Env *dotenv.Env
}

type UploadService struct {
	env   *dotenv.Env
	minio *minio.Client
}

func NewUploadService(opts UploadServiceOpts) *UploadService {
	s := &UploadService{
		env: opts.Env,
	}

	endpoint := s.env.MINIO_ENDPOINT
	accessKeyID := s.env.MINIO_ACCESS_KEY
	secretAccessKey := s.env.MINIO_SECRET_KEY

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil
	}
	s.minio = minioClient
	return s
}

func (u *UploadService) UploadFile(bucketName string, filename string, size int64, src multipart.File) error {
	exists, err := u.minio.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if !exists {
		err = u.minio.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	_, err = u.minio.PutObject(context.Background(), bucketName, filename, src, size, minio.PutObjectOptions{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
