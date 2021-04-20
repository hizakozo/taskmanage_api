package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
)

type ticketController struct {
	tu domain.TicketUsecase
}

func NewTicketController(tu domain.TicketUsecase) *ticketController {
	return &ticketController{
		tu: tu,
	}
}

func (tc *ticketController) Create (c echo.Context) error {
	user := interceptor.User
	form := &form.CreateTicketForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	worker := &form.Worker
	if form.Worker == 0 {
		worker = nil
	}
	tc.tu.CreateTicket(form.Title, form.Explanation, form.ProjectId, form.StatusId, &user.ID, worker)
	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func (tc *ticketController) GetList (c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	res, err := tc.tu.GetTicketList(projectId, user.ID, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (tc *ticketController) ChangeStatus(c echo.Context) error {
	user := interceptor.User
	form := &form.ChangeStatusForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	err := tc.tu.ChangeStatus(user.ID, form.ProjectId, form.StatusId, form.TicketId, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func (tc *ticketController) Update(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateTicketForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	res, err := tc.tu.Update(form.TicketId, user.ID, form.Worker, form.StatusId, form.Title, form.Explanation, c)
	if  err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (tc *ticketController) Detail (c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	res, err := tc.tu.Detail(ticketId, user.ID, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (tc *ticketController) Delete (c echo.Context) error {
	ticketId, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	err = tc.tu.Delete(ticketId, user.ID, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}