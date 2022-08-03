package auth_module

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"goauth/auth/auth_handler"
	"goauth/auth/auth_repository"
	"goauth/auth/auth_service"
	"goauth/user/user_repository"
	"goauth/user/user_service"
	"goauth/utils/datasources"
	"io/ioutil"
	"os"
)

type AuthModuleConfig struct {
	Group       *gin.RouterGroup
	DataSources *datasources.DataSources
}

func InitAuthModule(c *AuthModuleConfig) {
	userRepository := user_repository.NewUserRepository(c.DataSources.DB)
	userService := user_service.NewUserService(&user_service.ServiceConfig{
		UserRepository: userRepository,
	})
	passwordService := auth_handler.NewPassWordService()

	privateFile := os.Getenv("PRIVATE_KEY_FILE")
	private, _ := ioutil.ReadFile(privateFile)
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(private)

	publicFile := os.Getenv("PUBLIC_KEY_FILE")
	public, _ := ioutil.ReadFile(publicFile)
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(public)

	refreshSecret := os.Getenv("REFRESH_SECRET")

	tokenRepository := auth_repository.NewTokenRepository(c.DataSources.RedisClient)
	tokenService := auth_service.NewTokenService(&auth_service.TSConfig{
		PvKey:           privateKey,
		PubKey:          publicKey,
		RefreshSecret:   refreshSecret,
		TokenRepository: tokenRepository,
	})

	auth_handler.NewHandler(&auth_handler.Config{
		RouterGroup:     c.Group,
		UserService:     userService,
		PasswordService: passwordService,
		TokenService:    tokenService,
	})
}
