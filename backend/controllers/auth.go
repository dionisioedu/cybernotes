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

	log.Println("Senha digitada no login:", input.Password)
	log.Println("Hash salvo no banco:", user.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Println("Erro ao comparar senha:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha incorreta."})
		return
	}

	log.Println("Senha válida! Gerando token...")

	token := generateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(userID uint) string {
	claims := Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "cybernotes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("Erro ao gerar token:", err)
		return ""
	}

	return tokenString
}
