package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
)

type commentController struct {
	cu domain.CommentUsecase
}

func NewCommentController(cu domain.CommentUsecase) *commentController {
	return &commentController{
		cu: cu,
	}
}

func (cc *commentController) Detail(c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	res := cc.cu.Detail(ticketId)
	return c.JSON(http.StatusOK, res)
}

func (cc *commentController) Create(c echo.Context) error {
	req := &form.CreateCommentForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	res, err := cc.cu.Create(req.TicketId, user.ID, req.Comment, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
