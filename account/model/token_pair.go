package model

// TokenPair used for returning pairs of ids and refresh tokens.
type TokenPair struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}
