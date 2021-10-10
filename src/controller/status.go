package controller

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
)

type statusController struct {
	su domain.StatusUsecase
}

func NewStatusController(su domain.StatusUsecase) *statusController {
	return &statusController{
		su: su,
	}
}

func (sc *statusController) GetList(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	res, err := sc.su.GetList(user.ID, projectId, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (sc *statusController) Create(c echo.Context) error {
	user := interceptor.User
	req := &form.CreateStatusForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(req); err != nil {
		return exception.InputFailed(c)
	}
	res, err := sc.su.Create(user.ID, req.ProjectId, req.StatusName, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (sc *statusController) Update(c echo.Context) error {
	user := interceptor.User
	req := &form.UpdateStatusForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	res, err := sc.su.Update(user.ID, req.ProjectId, req.StatusId, req.Progress, req.StatusName, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (sc *statusController) Delete(c echo.Context) error {
	statusId, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		return exception.FormBindException(c)
	}

	user := interceptor.User
	err = sc.su.Delete(statusId, user.ID, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}