package models

type Users struct {
	UserId   string `gorm:"type:varchar(80);primary_key;"`
	Username string `gorm:"type:varchar(80);unique" form:"username" binding:"required"`
	Email    string `gorm:"type:varchar(80);unique" form:"email" binding:"required"`
	Password []byte
	RoleId   string `gorm:"type:varchar(80)"`
	Roles    Roles  `gorm:"foreignKey:RoleId"`
}

type Roles struct {
	IdRole string `gorm:"type:varchar(80);primary_key;"`
	Role   string `gorm:"type:varchar(80)"`
}
type AuthUserTokens struct {
	TokenId     string `gorm:"type:varchar(80);primary_key;"`
	IdUser      string `gorm:"type:varchar(80);unique" form:"id_user" binding:"required"`
	AccessToken string `gorm:"unique" binding:"required"`
	RefeshToken string `gorm:"unique" binding:"required"`
}