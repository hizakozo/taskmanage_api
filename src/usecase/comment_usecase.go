package usecase

import (
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/response"
)

type commentUsecase struct {
	ur domain.UserRepository
	tr domain.TicketRepository
	cr domain.CommentRepository
}

func NewCommentUsecase(ur domain.UserRepository, tr domain.TicketRepository, cr domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{
		ur: ur,
		tr: tr,
		cr: cr,
	}
}

func (cu *commentUsecase) Detail(ticketId int) response.CommentList {
	var responseComments []response.Comment
	findComments := cu.cr.CommentByTicketId(ticketId)
	for _, comment := range findComments {
		user, _ := cu.ur.UserById(comment.UserId)
		findCommentImgs := cu.cr.CommentImgByCommentId(comment.ID)
		var responseCommentImgs []response.CommentImg
		for _, commentImg := range findCommentImgs {
			responseCommentImg := response.CommentImg{Id: commentImg.ID, Path: commentImg.CommentImgPath}
			responseCommentImgs = append(responseCommentImgs, responseCommentImg)
		}
		responseUser := response.IdName{Id: user.ID, Name: user.Name}
		responseComment := response.Comment{Id: comment.ID, User: responseUser, Comment: comment.Comment, CommentImgs: responseCommentImgs}
		responseComments = append(responseComments, responseComment)
	}
	return response.CommentList{Comments: responseComments}
}

func (cu *commentUsecase) Create(ticketId, userId int, comment string, c echo.Context) (*response.CommentCreate, error) {
	err := cu.tr.TicketByIdUserId(ticketId, userId)
	if err != nil {
		return nil, exception.PermissionException(c)
	}
	_comment := domain.Comment{
		UserId:   userId,
		TicketId: ticketId,
		Comment:  comment}
	insertComment := cu.cr.InsertComment(_comment)
	return &response.CommentCreate{
		TicketId: insertComment.TicketId,
		Comment: response.Comment{Id: insertComment.ID,
			Comment: insertComment.Comment,
		},
	}, nil
}
