package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var uploader *s3manager.Uploader
var once sync.Once

func InitS3() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var initError error
	once.Do(func(){
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		})
		if err != nil {
			initError = fmt.Errorf("error creating session: %v", err)
			return
		}
		uploader = s3manager.NewUploader(sess)
	})
	return initError
}

func UpLoadImageToS3(file *multipart.FileHeader) (string, error) {

	if uploader == nil {
		if err := InitS3(); err != nil {
			return "", fmt.Errorf("error initializing s3: %v", err)
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	filename := fmt.Sprintf("%s%s", primitive.NewObjectID().Hex(), filepath.Ext(file.Filename))

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   src,
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	return result.Location, nil
}

func IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	return validTypes[contentType]
}
