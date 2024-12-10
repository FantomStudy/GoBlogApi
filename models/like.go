package models

type Like struct {
	Id     uint `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId uint `json:"user_id" gorm:"index"`
	PostId uint `json:"post_id" gorm:"index"`
	User   User `json:"user" gorm:"foreignKey:UserId"`
	Post   Post `json:"post" gorm:"foreignKey:PostId"`
}
