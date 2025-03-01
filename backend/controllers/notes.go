package controllers

import (
	"log"
	"net/http"

	"github.com/dionisioedu/cybernotes/backend/database"

	"github.com/gin-gonic/gin"
)

func CreateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	note := database.Note{
		UserID:  userID.(uint),
		Title:   input.Title,
		Content: input.Content,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		log.Println("Erro ao salvar nota:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar nota"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Nota criada com sucesso", "note": note})
}

func GetNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var notes []database.Note
	if err := database.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		log.Println("Erro ao buscar notas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar notas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func UpdateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	noteID := c.Param("id")
	var note database.Note
	if err := database.DB.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		log.Println("Nota não encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota não encontrada"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	note.Title = input.Title
	note.Content = input.Content

	if err := database.DB.Save(&note).Error; err != nil {
		log.Println("Erro ao atualizar nota:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar nota"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota atualizada com sucesso", "note": note})
}

func DeleteNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	noteID := c.Param("id")
	var note database.Note
	if err := database.DB.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		log.Println("Nota não encontrada:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota não encontrada"})
		return
	}

	if err := database.DB.Delete(&note).Error; err != nil {
		log.Println("Erro ao deletar nota:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar nota"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota deletada com sucesso"})
}
