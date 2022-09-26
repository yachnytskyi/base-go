package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// pgUserRepository is data/repository implementation of service layer UserRepository.
type pgUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository is a factory for initializating User Repositories.
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

// Create reaches out to database SQLX api.
func (r *pgUserRepository) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, user, query, user.Email, user.Password); err != nil {
		// Check unique constraint.
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, err.Code.Name())
			return apperrors.NewConflict("email", user.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByID fetches a user by id.
func (r *pgUserRepository) FindById(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE user_id=$1"

	// We need to actually check errors as it could be something other than not found.
	if err := r.DB.GetContext(ctx, user, query, userID); err != nil {
		return user, apperrors.NewNotFound("userID", userID.String())
	}

	return user, nil
}

// FindByEmail retrieves user row by email adrress.
func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get the user with email adress: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}
