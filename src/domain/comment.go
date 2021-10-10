package domain

import (
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/response"
)

type Comment struct {
	ID       int    `gorm:"column:comment_id;PRIMARY_KEY"`
	UserId   int    `gorm:"column:user_id"`
	TicketId int    `gorm:"column:ticket_id"`
	Comment  string `gorm:"column:comment"`
}

type CommentImg struct {
	ID             int    `gorm:"column:comment_img_id;PRIMARY_KEY"`
	CommentId      int    `gorm:"column:comment_id"`
	CommentImgPath string `gorm:"column:comment_img_path"`
}

type CommentRepository interface {
	CommentByTicketId(ticketId int) []Comment
	CommentImgByCommentId(commentId int) []CommentImg
	UpdateComment(commentId int, comment string)
	InsertComment(comment Comment) Comment
}

type CommentUsecase interface {
	Detail(ticketId int) response.CommentList
	Create(ticketId, userId int, comment string, c echo.Context) (*response.CommentCreate, error)
}