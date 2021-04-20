package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
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

type userController struct {
	uu domain.UserUseCase
}

func NewUserController(uu domain.UserUseCase) *userController {
	return &userController{
		uu: uu,
	}
}

func (uc *userController) SignUp(c echo.Context) error {
	req := form.SignUpForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error(err)
		return exception.InputFailed(c)
	}
	avatar, _ := c.FormFile("avatar")
	err := uc.uu.SignUp(
		c.FormValue("name"),
		c.FormValue("login_id"),
		c.FormValue("password"),
		c.FormValue("mail_address"),
		avatar,
		c,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func (uc *userController) SignIn(c echo.Context) error {
	req := &form.LoginForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(req); err != nil {
		return exception.InputFailed(c)
	}

	res, err := uc.uu.SignIn(req.LoginId, req.Password)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) GetUsersInProject(c echo.Context) error {
	user := interceptor.User
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return exception.FormBindException(c)
	}

	res, err := uc.uu.GetUsersInProject(user.ID, projectId, c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (uc *userController) SignOut(c echo.Context) error {
	token := c.Request().Header.Get("user_token")
	req := &form.LogOutForm{}
	if err := c.Bind(&req); err != nil {
		return exception.FormBindException(c)
	}
	err := uc.uu.SignOut(token, req.UserId, c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func (uc *userController) GetUserProfile(c echo.Context) error {
	userId := interceptor.User.ID
	res, err := uc.uu.GetUserProfile(userId, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
