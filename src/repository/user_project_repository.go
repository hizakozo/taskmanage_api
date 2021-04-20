package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type userProjectRepository struct {
	Db *gorm.DB
}

func NewUserProjectRepository(Db *gorm.DB) domain.UserProjectRepository {
	return &userProjectRepository{
		Db: Db,
	}
}

func (upr *userProjectRepository) InsertUserProject(userProject domain.UserProject) {
	upr.Db.Create(&userProject)
}

func (upr *userProjectRepository) UserProjectByUserIdProjectId(userId int, projectId int) []domain.UserProject {
	var userProject []domain.UserProject
	upr.Db.Select("user_project_id, user_id, project_id").
		Table("user_project").
		Where("user_id = ? AND project_id = ?", userId, projectId).
		Scan(&userProject)
	return userProject
}
