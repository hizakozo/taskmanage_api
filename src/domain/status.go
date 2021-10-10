package domain

import (
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/response"
)

type Status struct {
	ID         int    `gorm:"column:status_id;PRIMARY_KEY"`
	ProjectId  int    `gorm:"column:project_id"`
	Progress   int    `gorm:"column:progress"`
	StatusName string `gorm:"column:status_name"`
}

type StatusRepository interface {
	InsertStatus(status Status) Status
	StatusByProjectId(projectId int) []Status
	StatusById(statusId int) (*Status, error)
	StatusByTicketId(ticketId int) Status
	MaxProgressByProjectId(projectId int) int
	UpdateStatus(status Status) Status
	StatusByIdProjectId(statusId int, projectId int) (Status, error)
	UpdateProgress(statusId int, progress int)
	DeleteStatusTransaction(statusId int, statuses []Status, trgProgress int)
}

type StatusUsecase interface {
	GetList(userId, projectId int, c echo.Context) (*response.StatusList, error)
	Create(userId, projectId int, statusName string, c echo.Context) (*response.Status, error)
	Update(userId, projectId, statusId, progress int, statusName string, c echo.Context) (*response.Status, error)
	Delete(statusId, userId int, c echo.Context) error
}