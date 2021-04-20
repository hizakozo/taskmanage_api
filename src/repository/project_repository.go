package repository

import (
	"github.com/jinzhu/gorm"
	"taskmanage_api/src/domain"
)

type projectRepository struct {
	Db *gorm.DB
}

func NewProjectRepository(Db *gorm.DB) domain.ProjectRepository {
	return &projectRepository{
		Db: Db,
	}
}

func (pr *projectRepository) ProjectsByUserId(userId int) []domain.Project {
	var projects []domain.Project
	pr.Db.Table("project").Select("project.project_id, project_name, description, project_avatar").
		Joins("join user_project on project.project_id = user_project.project_id").
		Where("user_project.user_id = ?", userId).
		Scan(&projects)
	return projects
}

func (pr *projectRepository) ProjectById(projectId int) (*domain.Project, error) {
	project := domain.Project{ID: projectId}
	err := pr.Db.Find(&project).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &project, err
}

func (pr *projectRepository) InsertProject(project domain.Project) int {
	pr.Db.Create(&project)
	return project.ID
}

func (pr *projectRepository) UpdateProject(project domain.Project) {
	pr.Db.Save(&project)
}

func (pr *projectRepository) DeleteProject(projectId int) {
	project := domain.Project{ID: projectId}
	pr.Db.Delete(&project)
}