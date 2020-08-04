package handler

import (
	"encoding/json"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/data"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/form"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/mail"
	"taskmanage_api/src/model"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

func GetProjectList(c echo.Context) error {
	user := interceptor.User
	responseUser := response.User{Id: user.ID, Name: user.Name, Avatar: user.Avatar}

	var responseProjects []response.Project
	for _, project := range data.ProjectsByUserId(user.ID) {
		responseProjects =
			append(responseProjects, response.Project{Id: project.ID, Name: project.ProjectName,
				Description: project.Description, Avatar: project.ProjectAvatar})
	}
	return c.JSON(http.StatusOK,
		response.ProjectList{User: responseUser, Projects: responseProjects})
}

func CreateProject(c echo.Context) error {
	user := interceptor.User
	form := &form.CreateProjectForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(form); err != nil {
		return exception.InputFailed(c)
	}
	var project data.Project
	project.ProjectName = form.ProjectName
	project.Description = form.Description
	insertProjectId := data.InsertProject(project)

	var userProject data.UserProject
	userProject.UserId = user.ID
	userProject.ProjectId = insertProjectId
	data.InsertUserProject(userProject)

	statuses := constants.Statuses
	for _, v := range statuses {
		var status data.Status
		status.ProjectId = insertProjectId
		status.Progress = v.Progress
		status.StatusName = v.Name
		data.InsertStatus(status)
	}

	return c.JSON(http.StatusOK, response.SuccessResponse{Message: constants.ProcessingComplete})
}

func UpdateProject(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateProjectForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	if err := validator.New().Struct(form); err != nil {
		return exception.InputFailed(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}

	project := data.Project{ID: form.ProjectId, ProjectName: form.ProjectName,
		Description: form.Description, ProjectAvatar: form.ProjectAvatar}
	data.UpdateProject(project)

	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func DeleteProject(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if utils.IsErr(err) {
		return exception.FormBindException(c)
	}
	user := interceptor.User
	if userProject := data.UserProjectByUserIdProjectId(user.ID, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	data.DeleteProject(projectId)

	return c.JSON(http.StatusOK, constants.ProcessingComplete)
}

func InviteProject(c echo.Context) error {
	user := interceptor.User
	form := &form.InviteProjectForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}

	if userProject := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	auth, err := data.AuthByMailAddress(form.MailAddress)
	if utils.IsErr(err) {
		return exception.NotFoundData(c)
	}
	if userProject := data.UserProjectByUserIdProjectId(auth.ID, form.ProjectId); len(userProject) > 0 {
		return exception.DataAlreadyExists(c, "user_project")
	}
	inviteInfo := model.InviteInfo{ProjectId: form.ProjectId, UserId: auth.UserId}
	inviteInfoJson, _ := json.Marshal(inviteInfo)
	token, _ := utils.MakeRandomStr()
	data.RedisSet(string(inviteInfoJson), token)
	message :=
		constants.MailBody +
			constants.Params.FrontUrl + "#/join/" + token
	_ = mail.SendMail(auth.MailAddress, message)
	return c.JSON(http.StatusOK,
		response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})
}

func JoinProject(c echo.Context) error {
	form := &form.JoinProjectForm{}
	if err := c.Bind(&form); err != nil {
		return exception.FormBindException(c)
	}
	inviteInfoJson, _ := data.RedisGetInviteInfo(form.Token)

	var inviteInfo = new(model.InviteInfo)
	_ = json.Unmarshal([]byte(inviteInfoJson), inviteInfo)

	if userProject := data.UserProjectByUserIdProjectId(inviteInfo.UserId, inviteInfo.ProjectId); len(userProject) > 0 {
		return exception.DataAlreadyExists(c, "user_project")
	}

	data.InsertUserProject(data.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})

	data.RedisDelete(form.Token)

	return c.JSON(http.StatusOK,
		response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})
}
