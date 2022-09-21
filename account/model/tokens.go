package model

import "github.com/google/uuid"

// RefreshToken stores token properties that
// are accessed in multiple application layers.
type RefreshToken struct {
	ID           uuid.UUID `json:"-"`
	UID          uuid.UUID `json:"-"`
	SignedString string    `json:"refreshToken"`
}

// IDToken stores token properties that
// are accessed in multiple application layers.
type IDToken struct {
	SignedString string `json:"idToken"`
}

// TokenPair used for returning pairs of id and refresh tokens.
type TokenPair struct {
	IDToken
	RefreshToken
}
