package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

type googleCloudImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

// NewImageRepository is a factory for initializing User Repositories.
func NewImageRepository(googleCloudClient *storage.Client, bucketName string) model.ImageRepository {
	return &googleCloudImageRepository{
		Storage:    googleCloudClient,
		BucketName: bucketName,
	}
}

func (r *googleCloudImageRepository) UpdateProfile(ctx context.Context, objectName string, imageFile multipart.File) (string, error) {
	bucket := r.Storage.Bucket(r.BucketName)

	object := bucket.Object(objectName)
	writerStorage := object.NewWriter(ctx)

	// Set cache control so a profile image will be served fresh by browsers.
	// To do this with an object handle, we would first have to upload, then update.
	writerStorage.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0"

	// multipart.File has a reader!
	if _, err := io.Copy(writerStorage, imageFile); err != nil {
		log.Printf("Unable to write file to Google Cloud Storage: %v\n", err)
		return "", apperrors.NewInternal()
	}

	if err := writerStorage.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", r.BucketName, objectName)

	return imageURL, nil
}
