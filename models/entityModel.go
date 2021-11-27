package models

type Users struct {
	UserId   string `gorm:"type:varchar(80);primary_key;"`
	Username string `gorm:"type:varchar(80);unique" form:"username" binding:"required"`
	Email    string `gorm:"type:varchar(80);unique" form:"email" binding:"required"`
	Password []byte
}

type AuthUserTokens struct {
	TokenId     string `gorm:"type:varchar(80);primary_key;"`
	UserId      string `gorm:"type:varchar(80);unique" form:"user_id" binding:"required"`
	AccessToken string `gorm:"unique" binding:"required"`
	RefeshToken string `gorm:"unique" binding:"required"`
}