package data

import "github.com/jinzhu/gorm"

type Auth struct {
	ID          int    `gorm:"column:auth_id;PRIMARY_KEY"`
	UserId      int    `gorm:"column:user_id"`
	LoginId     string `gorm:"column:login_id"`
	Password    string
	MailAddress string
}

func InsertAuth(auth Auth) {
	Db.Create(&auth)
}

func AuthByLoginId(loginId string) (*Auth, error){
	auth := Auth{}
	err := Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("login_id = ?", loginId).
		Find(&auth).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &auth, err
}

func AuthByUserId(userId int) (Auth, error) {
	auth := Auth{}
	err := Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("user_id = ?", userId).
		Find(&auth).Error
	return auth, err
}

func AuthByMailAddress(mailAddress string) (*Auth, error) {
	auth := Auth{}
	err := Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("mail_address = ?", mailAddress).
		Find(&auth).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &auth, err
}