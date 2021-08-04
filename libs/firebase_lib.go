package libs

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirebaseService interface {
	UploadFile(*multipart.FileHeader, string) error
}

type FirebaseLib struct {
	firebaseService FirebaseService
}

func (f FirebaseLib) BeginUpload(fileHeader *multipart.FileHeader, bucketName string) error {
	if err := f.firebaseService.UploadFile(fileHeader, bucketName); err != nil {
		return err
	}

	return nil
}

var UploadLib *FirebaseLib

// Implementation of FirebaseService - UploadFile
type UploadService struct {
	App *firebase.App
}

func (u UploadService) UploadFile(fileHeader *multipart.FileHeader, bucketName string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	ctx := context.Background()
	storage, err := u.App.Storage(ctx)
	if err != nil {
		return err
	}

	bkt, err := storage.Bucket(os.Getenv("GCS_BUCKET_NAME"))
	if err != nil {
		return err
	}

	obj := bkt.Object(bucketName)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func ConnectFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("firebaseServiceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// can be called from unit test to mock the upload behavior
func InitUploadLib(service FirebaseService) {
	UploadLib = &FirebaseLib{firebaseService: service}
}
