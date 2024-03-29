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

func (repository *googleCloudImageRepository) DeleteProfile(ctx context.Context, objectName string) error {
	bucket := repository.Storage.Bucket(repository.BucketName)

	object := bucket.Object(objectName)

	if err := object.Delete(ctx); err != nil {
		log.Printf("Failed to delete the image object with ID: %s from Google Cloud Storage\n", objectName)
		return apperrors.NewInternal()
	}

	return nil
}

func (repository *googleCloudImageRepository) UpdateProfile(ctx context.Context, objectName string, imageFile multipart.File) (string, error) {
	bucket := repository.Storage.Bucket(repository.BucketName)

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

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", repository.BucketName, objectName)

	return imageURL, nil
}
