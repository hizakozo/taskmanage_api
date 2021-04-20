package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type userRepository struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) domain.UserRepository {
	return &userRepository{
		Db: Db,
	}
}

func(r *userRepository) UserById(userId int) (domain.User, error) {
	user := domain.User{}
	err := r.Db.Select("user_id, user_name, avatar").
		Table("user").
		Where("user_id  = ?", userId).
		Find(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return domain.User{}, err
	}
	return user, nil
}

func(r *userRepository)  InsertUser(user domain.User) int {
	r.Db.Create(&user)
	return user.ID
}

func(r *userRepository)  UserByProjectId(projectId int) []domain.User {
	var users []domain.User
	r.Db.Select("u.user_id, user_name, avatar").
		Table("user u").
		Joins("join user_project up on u.user_id = up.user_id").
		Where("up.project_id = ?", projectId).
		Scan(&users)
	return users
}