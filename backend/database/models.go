package database

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token  string `json:"token" gorm:"unique"`
	UserID uint   `json:"user_id"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"text"`
}

type Note struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content" gorm:"type:text"`
}
