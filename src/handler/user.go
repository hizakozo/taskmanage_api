package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/data"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

func SignUp(c echo.Context) error {
	form := form.SignUpForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(form); err != nil {
		log.Error(err)
		return exception.InputFailed(c)
	}
	name := c.FormValue("name")
	loginId := c.FormValue("login_id")
	password := c.FormValue("password")
	mailAddress := c.FormValue("mail_address")
	avatar, _ := c.FormFile("avatar")
	if _, err := data.AuthByLoginId(loginId); err == nil {
		return exception.DataAlreadyExists(c, "login id")
	}
	if _, err := data.AuthByMailAddress(mailAddress); err == nil {
		return exception.DataAlreadyExists(c, "mail address")
	}
	user := data.User{Name: name}
	if avatar != nil {
		fileName := "user/" + avatar.Filename
		src, err := avatar.Open()
		defer src.Close()
		_, err = utils.S3PutObject(fileName, src)
		if err != nil {
			fmt.Print(err)
			return exception.FileUploadFailed(c)
		}
		user.Avatar = fileName
	}
	id := data.InsertUser(user)
	safetyPass := utils.CreateSafetyPass(password)
	auth := data.Auth{UserId: id, LoginId: loginId, Password: safetyPass, MailAddress: mailAddress}

	data.InsertAuth(auth)

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func SignIn(c echo.Context) error {
	form := &form.LoginForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(form); err != nil {
		return exception.InputFailed(c)
	}
	auth, err := data.AuthByLoginId(form.LoginId)
	fmt.Print(form)
	err = utils.PasswordVerify(auth.Password, form.Password)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	user, err := data.UserById(auth.UserId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	var userToken, _ = utils.MakeRandomStr()
	userJson, _ := json.Marshal(user)
	data.RedisSet(string(userJson), userToken)

	return c.JSON(http.StatusOK, response.LoginResponse{UserToken: userToken, UserId: user.ID})
}

func GetUsersInProject(c echo.Context) error {
	user := interceptor.User
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	users := data.UserByProjectId(projectId)
	var responseUsers []response.User
	for _, user := range users {
		responseUser := response.User{Id: user.ID, Name: user.Name}
		if user.Avatar != "" {
			responseUser.Avatar = constants.Params.S3Url + user.Avatar
		}
		responseUsers = append(responseUsers, responseUser)
	}

	return c.JSON(http.StatusOK, response.UserList{Users: responseUsers})
}

func SignOut(c echo.Context) error {
	token := c.Request().Header.Get("user_token")
	user, _ := data.RedisGet(token)
	form := &form.LogOutForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if user.ID != form.UserId {
		return exception.PermissionException(c)
	}
	data.RedisDelete(token)
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func GetUserProfile(c echo.Context) error {
	userId := interceptor.User.ID
	user, err := data.UserById(userId)
	auth, err := data.AuthByUserId(userId)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	projects := data.ProjectsByUserId(userId)
	responseUser := response.User{Id: user.ID, Name: user.Name,
		Avatar: user.Avatar, MailAddress: auth.MailAddress}

	var responseProjects []response.ProjectTickets
	for _, project := range projects {
		var responseTickets []response.Ticket
		for _, ticket := range data.TicketByProjectIdWorker(project.ID, user.ID) {
			responseTickets = append(responseTickets, response.Ticket{Id: ticket.ID, Title: ticket.Title})
		}
		responseProjects = append(responseProjects, response.ProjectTickets{Id: project.ID, Name: project.ProjectName,
			Tickets: responseTickets})
	}

	return c.JSON(http.StatusOK, response.UserProfile{User: responseUser, Projects: responseProjects})
}
