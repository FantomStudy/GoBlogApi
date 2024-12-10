package models

type Photo struct {
	Id     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Path   string `json:"path"`
	PostId uint   `json:"post_id" gorm:"index"`
	Post   Post   `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}
