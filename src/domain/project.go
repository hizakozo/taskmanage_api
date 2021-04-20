package domain

import (
	"github.com/labstack/echo"
	"taskmanage_api/src/response"
)

type Project struct {
	ID            int    `gorm:"column:project_id;PRIMARY_KEY"`
	ProjectName   string `gorm:"column:project_name"`
	Description   string `gorm:"column:description"`
	ProjectAvatar string `gorm:"column:project_avatar"`
}

type ProjectRepository interface {
	ProjectsByUserId(userId int) []Project
	ProjectById(projectId int) (*Project, error)
	InsertProject(project Project) int
	UpdateProject(project Project)
	DeleteProject(projectId int)
}

type ProjectUsecase interface {
	GetList(userId int, userName, avatar string) *response.ProjectList
	Create(projectName, description string, userId int) error
	Update(userId, projectId int, projectName, description, avatar string, c echo.Context) error
	Delete(userId, projectId int, c echo.Context) error
	Invite(userId, projectId int, mailAddress string, c echo.Context) (*response.UserProject, error)
	Join(token string, c echo.Context) (*response.UserProject, error)
}
