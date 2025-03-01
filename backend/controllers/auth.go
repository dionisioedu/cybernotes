package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dionisioedu/cybernotes/backend/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func Register(c *gin.Context) {
	var user database.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	log.Println("Senha recebida no registro:", user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash da senha:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criptografar senha"})
		return
	}

	log.Println("Hash gerado:", string(hashedPassword))
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("Erro ao salvar usuário:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

func Login(c *gin.Context) {
	var user database.User
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	input.Password = strings.TrimSpace(input.Password)

	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Println("Usuário não encontrado:", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Println("Erro ao comparar senha:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha incorreta."})
		return
	}

	accessToken := generateToken(user.ID, 15)
	refreshToken := generateToken(user.ID, 60*24*7) // Expires in 7 days

	database.DB.Create(&database.RefreshToken{
		Token:  refreshToken,
		UserID: user.ID,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func generateToken(userID uint, minutes int) string {
	claims := Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(minutes)).Unix(),
			Issuer:    "cybernotes",
		},
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY não configurada")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Erro ao gerar token:", err)
		return ""
	}

	return tokenString
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var token database.RefreshToken
	if err := database.DB.Where("token = ?", input.RefreshToken).First(&token).Error; err != nil {
		log.Println("Token inválido:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	accessToken := generateToken(token.UserID, 15)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
