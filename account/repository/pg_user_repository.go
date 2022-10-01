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

// pgUserRepository is data/repository implementation of the service layer UserRepository.
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
func (r *pgUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
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

// Update updates a user's properties.
func (r *pgUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET username=:username, email=:email, website=:website
		WHERE user_id=:user_id
		RETURNING *;
	`
	prepareNamedStatement, err := r.DB.PrepareNamedContext(ctx, query)

	if err != nil {
		log.Printf("Unable to prepate the user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := prepareNamedStatement.GetContext(ctx, user, user); err != nil {
		log.Printf("Unable to prepare the user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	return nil
}

// UpdateImage is used to update a user's separately image from
// other account details.
func (r *pgUserRepository) UpdateImage(ctx context.Context, userID uuid.UUID, imageURL string) (*model.User, error) {
	query := `
		UPDATE users
		SET image_url=$2
		WHERE user_id=$1
		RETURNING *;
	`

	// Must be instantiated to scan into reference using 'GetContext'.
	user := &model.User{}

	err := r.DB.GetContext(ctx, user, query, userID, imageURL)

	if err != nil {
		log.Printf("Error updating image_url in database: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}
