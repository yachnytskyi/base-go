package repository

import (
	"cloud.google.com/go/storage"
	"github.com/yachnytskyi/base-go/account/model"
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
