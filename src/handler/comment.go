package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"taskmanage_api/src/data"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

func GetComment(c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}
	var responseComments []response.Comment
	findComments := data.CommentByTicketId(ticketId)
	for _, comment := range findComments {
		user, _ := data.UserById(comment.UserId)
		findCommentImgs := data.CommentImgByCommentId(comment.ID)
		var responseCommentImgs []response.CommentImg
		for _, commentImg := range findCommentImgs {
			responseCommentImg := response.CommentImg{Id: commentImg.ID, Path: commentImg.CommentImgPath}
			responseCommentImgs = append(responseCommentImgs, responseCommentImg)
		}
		responseUser := response.IdName{Id: user.ID, Name: user.Name}
		responseComment := response.Comment{Id: comment.ID, User: responseUser, Comment: comment.Comment, CommentImgs: responseCommentImgs}
		responseComments = append(responseComments, responseComment)
	}

	return c.JSON(http.StatusOK, response.CommentList{Comments: responseComments})
}

func CreateComment(c echo.Context) error {
	form := &form.CreateCommentForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	err := data.TicketByIdUserId(form.TicketId, user.ID)
	if utils.IsErr(err) {
		return exception.PermissionException(c)
	}
	comment := data.Comment{UserId: user.ID, TicketId: form.TicketId, Comment: form.Comment}
	insertComment := data.InsertComment(comment)
	return c.JSON(http.StatusOK,
		response.CommentCreate{TicketId: insertComment.TicketId, Comment:
		response.Comment{Id: insertComment.ID, Comment: insertComment.Comment}})
}
