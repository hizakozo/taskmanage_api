package handler

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"taskmanage_api/constants"
	"taskmanage_api/data"
	"taskmanage_api/form"
	"taskmanage_api/interceptor"
	"taskmanage_api/mail"
	"taskmanage_api/model"
	"taskmanage_api/response"
	"taskmanage_api/utils"
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
	if err := c.Bind(form); err != nil {
		return err
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

	return c.JSON(http.StatusOK, "create project")
}

func UpdateProject(c echo.Context) error {
	user := interceptor.User
	form := &form.UpdateProjectForm{}
	if err := c.Bind(form); err != nil {
		return err
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}

	project := data.Project{ID: form.ProjectId, ProjectName: form.ProjectName,
		Description: form.Description, ProjectAvatar: form.ProjectAvatar}
	data.UpdateProject(project)

	return c.JSON(http.StatusOK, "project update")
}

func DeleteProject(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	user := interceptor.User
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	data.DeleteProject(projectId)

	return c.JSON(http.StatusOK, "ticket delete")
}

func InviteProject(c echo.Context) error {
	user := interceptor.User
	form := &form.InviteProjectForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: constants.PermissionException})
	}
	auth, err := data.AuthByMailAddress(form.MailAddress)
	if utils.IsErr(err) {
		return response.CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(auth.ID, form.ProjectId); err == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "already exists user_project"})
	}
	inviteInfo := model.InviteInfo{ProjectId: form.ProjectId, UserId: auth.UserId}
	inviteInfoJson, err := json.Marshal(inviteInfo)
	if err != nil {
		return response.CreateErrorResponse(err, c)
	}
	token, _ := utils.MakeRandomStr()
	data.RedisSet(string(inviteInfoJson), token)
	message := "http://localhost:8081/join/" + token
	if err := mail.SendMail(auth.MailAddress, message); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, response.SuccessResponse{Message: "send mail to" + auth.MailAddress})
}

func JoinProject(c echo.Context) error {
	form := &form.JoinProjectForm{}
	if err := c.Bind(form); err != nil {
		return response.CreateErrorResponse(err, c)
	}
	inviteInfoJson, _ := data.RedisGetInviteInfo(form.Token)

	var inviteInfo = new(model.InviteInfo)
	if err := json.Unmarshal([]byte(inviteInfoJson), inviteInfo); err != nil {
		return response.CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(inviteInfo.UserId, inviteInfo.ProjectId); err == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "already exists user_project"})
	}

	data.InsertUserProject(data.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})

	return c.JSON(http.StatusOK,
		response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})
}
