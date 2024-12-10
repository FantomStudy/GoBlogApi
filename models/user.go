package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	Phone    string `json:"phone"`
}

func (user *User) SetPassword(password string) {
	// Хэширование пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	// Сравнение пароля из базы с введенным
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
