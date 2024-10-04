package login

import (
	"time"
	models_user "tm/controllers/admin/login/models"
	config "tm/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")
var refreshSecret = []byte("your-refresh-secret")

func generateTokens(userID uint) (accessToken string, refreshToken string, err error) {
	accessTokenClaims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err = accessJwt.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken, err = refreshJwt.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var user models_user.User
	config.DB.Where("username = ?", data["username"]).First(&user)
	if user.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	newUser := models_user.User{
		Username: data["username"],
		Password: data["password"],
	}

	config.DB.Create(&newUser)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var user models_user.User
	config.DB.Where("username = ? AND password = ?", data["username"], data["password"]).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Logged in successfully",
	})
}

func Protected(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or invalid access token",
		})
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid access token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	c.Locals("userID", uint(userID))

	return c.Next()
}

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or invalid refresh token",
		})
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	// Kullanıcı ID'sini refresh token'dan al
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	// Yeni access token oluştur
	newAccessToken, _, err := generateTokens(uint(userID))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Yeni access token cookie'ye ekleniyor
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Access token refreshed",
	})
}
