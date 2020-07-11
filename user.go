package main

import (
	"./constants"
	"./data"
	"./interceptor"
	"./response"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type SignUpForm struct {
	Name        string `json:"name"`
	LoginId     string `json:"login_id"`
	Password    string `json:"password"`
	MailAddress string `json:"mail_address"`
}

type LoginForm struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

func signUp(c echo.Context) error {
	form := &SignUpForm{}

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
	safetyPass := createSafetyPass(form.Password)
	auth.Password = safetyPass
	auth.MailAddress = form.MailAddress

	data.InsertAuth(auth)

	return c.JSON(http.StatusOK, id)
}

func login(c echo.Context) error {
	form := &LoginForm{}
	if err := c.Bind(form); err != nil {
		return err
	}
	auth, err := data.AuthByLoginId(form.LoginId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if err := passwordVerify(auth.Password, form.Password); err != nil {
		return CreateErrorResponse(err, c)
	}

	user, err := data.UserById(auth.UserId)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	var userToken, _ = MakeRandomStr()

	userJson, _ := json.Marshal(user)

	data.RedisSet(string(userJson), userToken)

	return c.JSON(http.StatusOK, response.LoginResponse{UserToken: userToken})
}

func getUsersInProject(c echo.Context) error {
	user := interceptor.User
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	users := data.UserByProjectId(projectId)
	var responseUsers []response.IdName
	for _, user := range users {
		responseUser := response.IdName{Id: user.ID, Name: user.Name}
		responseUsers = append(responseUsers, responseUser)
	}

	return c.JSON(http.StatusOK, response.UserList{Users: responseUsers})
}
