package models

type Post struct {
	Id     uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string  `json:"title" gorm:"index"`
	Desc   string  `json:"desc"`
	UserId string  `json:"user_id" gorm:"index"`
	User   User    `json:"user" gorm:"foreignKey:UserId"`
	Photos []Photo `json:"photos" gorm:"foreignKey:PostId"`
}
