package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/handler"
	"github.com/yachnytskyi/base-go/account/repository"
	"github.com/yachnytskyi/base-go/account/service"
)

// Will initialize a handler starting from data sources
// which inject into the repository layer
// which inject into the service layer
// which inject into the handler layer.
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injection data sources")

	/*
	 * repository layer.
	 */
	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)

	bucketName := os.Getenv("GOOGLE_CLOUD_IMAGE_BUCKET")
	imageRepository := repository.NewImageRepository(d.StorageClient, bucketName)

	/*
	 * service layer.
	 */
	userService := service.NewUserService(&service.UserConfig{
		UserRepository:  userRepository,
		ImageRepository: imageRepository,
	})

	// Load rsa keys.
	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	private, err := ioutil.ReadFile(privateKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(private)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	publicKeyFile := os.Getenv("PUBLIC_KEY_FILE")
	public, err := ioutil.ReadFile(publicKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(public)

	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// Load refresh token secret from env variable.
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// Load expiration lengts from env variables and parse as int.
	idTokenExpiration := os.Getenv("ID_TOKEN_EXPIRATION")
	refreshTokenExpiration := os.Getenv("REFRESH_TOKEN_EXPIRATION")

	idExpiration, err := strconv.ParseInt(idTokenExpiration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXPIRATION as int: %w", err)
	}

	refreshExpiration, err := strconv.ParseInt(refreshTokenExpiration, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXPIRATION as int: %w", err)
	}

	tokenService := service.NewTokenService(&service.TokenServiceConfig{
		TokenRepository:          tokenRepository,
		PrivateKey:               privateKey,
		PublicKey:                publicKey,
		RefreshSecret:            refreshSecret,
		IDExpirationSecrets:      idExpiration,
		RefreshExpirationSecrets: refreshExpiration,
	})

	// Initialize gin.Engine
	router := gin.Default()

	// Read in ACCOUNT_API_URL.
	baseURL := os.Getenv("ACCOUNT_API_URL")

	// Read in HANDLER_TIMEOUT.
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	handlerTimeoutInt, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	maxBodyBytesParsed, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		Router:          router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(time.Duration(handlerTimeoutInt) * time.Second),
		MaxBodyBytes:    maxBodyBytesParsed,
	})

	return router, nil

}
