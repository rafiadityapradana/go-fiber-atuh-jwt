package models

type Users struct {
	UserId         string `gorm:"type:varchar(80);primary_key;"`
	Username       string `gorm:"type:varchar(80);unique" form:"username" binding:"required"`
	Email          string `gorm:"type:varchar(80);unique" form:"email" binding:"required"`
	Password       []byte
	RoleId         string         `gorm:"type:varchar(80)"`
	Roles          Roles          `gorm:"foreignKey:RoleId"`
	AuthUserTokens AuthUserTokens `gorm:"foreignKey:TokenUserId"`
}

type Roles struct {
	IdRole string `gorm:"type:varchar(80);primary_key;"`
	Role   string `gorm:"type:varchar(80)"`
}
type AuthUserTokens struct {
	TokenUserId string `gorm:"type:varchar(80);primary_key"`
	AccessToken string `gorm:"unique" binding:"required"`
	RefeshToken string `gorm:"unique" binding:"required"`
}