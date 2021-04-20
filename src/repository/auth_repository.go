package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type authRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) *authRepository {
	return &authRepository{
		Db: Db,
	}
}

func (ar *authRepository) InsertAuth(auth domain.Auth) {
	ar.Db.Create(&auth)
}

func (ar *authRepository) AuthByLoginId(loginId string) (domain.Auth, error) {
	auth := domain.Auth{}
	err := ar.Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("login_id = ?", loginId).
		Find(&auth).Error
	if gorm.IsRecordNotFoundError(err) {
		return domain.Auth{}, err
	}
	return auth, nil
}

func (ar *authRepository) AuthByUserId(userId int) (domain.Auth, error) {
	auth := domain.Auth{}
	err := ar.Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("user_id = ?", userId).
		Find(&auth).Error
	if gorm.IsRecordNotFoundError(err) {
		return domain.Auth{}, err
	}
	return auth, nil
}

func (ar *authRepository) AuthByMailAddress(mailAddress string) (domain.Auth, error) {
	auth := domain.Auth{}
	err := ar.Db.Select("auth_id, user_id, login_id, password, mail_address").
		Table("auth").
		Where("mail_address = ?", mailAddress).
		Find(&auth).Error
	if gorm.IsRecordNotFoundError(err) {
		return domain.Auth{}, err
	}
	return auth, nil
}
