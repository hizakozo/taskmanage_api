package handler

import (
	"taskmanage_api/constants"
	"taskmanage_api/data"
	"taskmanage_api/form"
	"taskmanage_api/interceptor"
	"taskmanage_api/response"
	"taskmanage_api/utils"
	"encoding/json"
	"github.com/labstack/echo"
	"mime/multipart"
	"net/http"
	"strconv"
)

func SignUp(c echo.Context) error {
	form := &form.SignUpForm{}

	if err := c.Bind(form); err != nil {
		return err
	}
	var user data.User
	user.Name = form.Name
	user.Avatar = "avatar"

	id := data.InsertUser(user)

	var auth data.Auth
	auth.UserId = id
	auth.LoginId = form.LoginId
	safetyPass := utils.CreateSafetyPass(form.Password)
	auth.Password = safetyPass
	auth.MailAddress = form.MailAddress

	data.InsertAuth(auth)

	return c.JSON(http.StatusOK, id)
}

func SignIn(c echo.Context) error {
	form := &form.LoginForm{}
	if err := c.Bind(form); err != nil {
		return err
	}
	auth, err := data.AuthByLoginId(form.LoginId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if err := utils.PasswordVerify(auth.Password, form.Password); err != nil {
		return response.CreateErrorResponse(err, c)
	}

	user, err := data.UserById(auth.UserId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	var userToken, _ = utils.MakeRandomStr()

	userJson, _ := json.Marshal(user)

	data.RedisSet(string(userJson), userToken)

	return c.JSON(http.StatusOK, response.LoginResponse{UserToken: userToken, UserId: user.ID})
}

func GetUsersInProject(c echo.Context) error {
	user := interceptor.User
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{constants.PermissionException})
	}
	users := data.UserByProjectId(projectId)
	var responseUsers []response.IdName
	for _, user := range users {
		responseUser := response.IdName{Id: user.ID, Name: user.Name}
		responseUsers = append(responseUsers, responseUser)
	}

	return c.JSON(http.StatusOK, response.UserList{Users: responseUsers})
}

func SignOut(c echo.Context) error {
	token := c.Request().Header.Get("user_token")
	user, err := data.RedisGet(token)
	form := &form.LogOutForm{}
	err = c.Bind(form)
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	if user.ID != form.UserId {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "bad Request"})
	}
	data.RedisDelete(token)
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: "sign out"})
}

type Up struct {
	File multipart.File
}