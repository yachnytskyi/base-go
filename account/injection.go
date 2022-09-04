package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/handler"
	"github.com/yachnytskyi/base-go/account/repository"
	"github.com/yachnytskyi/base-go/account/service"
)

// Will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer.
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injection data sources")

	/*
	 * repository layer.
	 */
	userRepository := repository.NewUserRepository(d.DB)

	/*
	 * repository layer.
	 */
	userService := service.NewUserService(&service.UserConfig{
		UserRepository: userRepository,
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

	tokenService := service.NewTokenService(&service.TokenServiceConfig{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		RefreshSecret: refreshSecret,
	})

	// Initialize gin.Engine
	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R:            router,
		UserService:  userService,
		TokenService: tokenService,
	})

	return router, nil

}
