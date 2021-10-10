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

type projectController struct {
	pu domain.ProjectUsecase
}

func NewProjectController(pu domain.ProjectUsecase) *projectController {
	return &projectController{
		pu: pu,
	}
}

func (pc *projectController) GetList(c echo.Context) error {
	user := interceptor.User
	res := pc.pu.GetList(user.ID, user.Name, user.Avatar)
	return c.JSON(http.StatusOK, res)
}

func (pc *projectController) Create(c echo.Context) error {
	user := interceptor.User
	req := &form.CreateProjectForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(req); err != nil {
		return exception.InputFailed(c)
	}
	err := pc.pu.Create(req.ProjectName, req.Description, user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func (pc *projectController) Update(c echo.Context) error {
	user := interceptor.User
	req := &form.UpdateProjectForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(req); err != nil {
		return exception.InputFailed(c)
	}
	err := pc.pu.Update(user.ID, req.ProjectId, req.ProjectName, req.Description, req.ProjectAvatar, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func (pc *projectController) Delete(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	err = pc.pu.Delete(user.ID, projectId, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func (pc *projectController) Invite(c echo.Context) error {
	user := interceptor.User
	req := &form.InviteProjectForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}

	res, err := pc.pu.Invite(user.ID, req.ProjectId, req.MailAddress, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (pc *projectController) Join(c echo.Context) error {
	req := &form.JoinProjectForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	res, err := pc.pu.Join(req.Token, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
