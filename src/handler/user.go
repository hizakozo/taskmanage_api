package handler

import (
	"encoding/json"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/data"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

func SignUp(c echo.Context) error {
	form := &form.SignUpForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	if err := validator.New().Struct(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	name := c.FormValue("name")
	loginId := c.FormValue("login_id")
	password := c.FormValue("password")
	mailAddress := c.FormValue("mail_address")
	avatar, _ := c.FormFile("avatar")
	if _, err := data.AuthByLoginId(loginId); err == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "login id is already exists"})
	}
	if _, err := data.AuthByMailAddress(mailAddress); err == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "mail address id is already exists"})
	}
	user :=  data.User{Name: name}
	if avatar != nil {
		fileName := "/user/" + avatar.Filename
		src, err := avatar.Open()
		defer src.Close()
		_, err = utils.S3PutObject("taskmanage", fileName, src)
		if err != nil {
			return response.CreateErrorResponse(err, c)
		}
		user.Avatar = fileName
	}
	id := data.InsertUser(user)
	safetyPass := utils.CreateSafetyPass(password)
	auth := data.Auth{UserId: id, LoginId: loginId, Password: safetyPass, MailAddress: mailAddress}

	data.InsertAuth(auth)

	return c.JSON(http.StatusOK, id)
}

func SignIn(c echo.Context) error {
	form := &form.LoginForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	if err := validator.New().Struct(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	auth, err := data.AuthByLoginId(form.LoginId)
	err = utils.PasswordVerify(auth.Password, form.Password)
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
	var responseUsers []response.User
	for _, user := range users {
		responseUser := response.User{Id: user.ID, Name: user.Name}
		if user.Avatar != "" {
			responseUser.Avatar = "http://127.0.0.1:4572/taskmanage" + user.Avatar
		}
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

func GetUserProfile(c echo.Context) error {
	userId := interceptor.User.ID
	user, err := data.UserById(userId)
	auth, err := data.AuthByUserId(userId)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
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