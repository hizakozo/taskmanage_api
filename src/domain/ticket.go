package domain

import (
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/response"
)

type Ticket struct {
	ID          int    `gorm:"column:ticket_id;PRIMARY_KEY"`
	ProjectId   int    `gorm:"column:project_id"`
	Title       string `gorm:"column:title"`
	Explanation string `gorm:"column:explanation"`
	Reporter    *int   `gorm:"column:reporter"`
	Worker      *int   `gorm:"column:worker"`
}

type TicketImg struct {
	ID            int    `gorm:"column:ticket_img_id;PRIMARY_KEY"`
	TicketId      int    `gorm:"column:ticket_id"`
	TicketImgPath string `gorm:"column:ticket_img_path"`
}

type TicketRepository interface {
	InsertTicket(ticket Ticket) int
	TicketByProjectIdStatusId(projectId int, statusId int) []Ticket
	TicketById(ticketId int) (*Ticket, error)
	TicketByProjectIdWorker(projectId int, worker int) []Ticket
	UpdateTicket(ticket *Ticket, ticketStatusId int, statusId int)
	TicketImgById(ticketId int) []TicketImg
	DeleteTicket(ticketId int)
	TicketByStatusId(statusId int) []Ticket
	TicketByIdUserId(ticketId int, userId int) error
}

type TicketUsecase interface {
	CreateTicket (title, explanation string, projectId, statusId int, reporter, worker *int)
	GetTicketList(projectId, userId int, c echo.Context) (*response.TicketList, error)
	ChangeStatus(userId, projectId, statusId, ticketId int, c echo.Context) error
	Update(ticketId, userId, worker, statusId int, title, explanation string, c echo.Context) (*response.Ticket, error)
	Detail(ticketId, userId int, c echo.Context) (*response.TicketDetail, error)
	Delete (ticketId, userId int, c echo.Context) error
}
