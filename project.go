package main

import (
	"./constants"
	"./data"
	"./interceptor"
	"./mail"
	"./response"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type CreateProjectForm struct {
	ProjectName string `json:"project_name"`
	Description string `json:"description"`
}

type UpdateProjectForm struct {
	ProjectId     int    `json:"project_id"`
	ProjectName   string `json:"project_name"`
	Description   string `json:"description"`
	ProjectAvatar string `json:"project_avatar"`
}

type InviteProjectForm struct {
	ProjectId   int    `json:"project_id"`
	MailAddress string `json:"mail_address"`
}

type JoinProjectForm struct {
	Token string `json:"token"`
}

type InviteInfo struct {
	ProjectId int `json:"project_id"`
	UserId    int `json:"user_id"`
}

func getProjectList(c echo.Context) error {
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

func createProject(c echo.Context) error {
	user := interceptor.User
	form := &CreateProjectForm{}
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

func updateProject(c echo.Context) error {
	user := interceptor.User
	form := &UpdateProjectForm{}
	if err := c.Bind(form); err != nil {
		return err
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}

	project := data.Project{ID: form.ProjectId, ProjectName: form.ProjectName,
		Description: form.Description, ProjectAvatar: form.ProjectAvatar}
	data.UpdateProject(project)

	return c.JSON(http.StatusOK, "project update")
}

func deleteProject(c echo.Context) error {
	projectId, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	user := interceptor.User
	if err := data.UserProjectByUserIdProjectId(user.ID, projectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	data.DeleteProject(projectId)

	return c.JSON(http.StatusOK, "ticket delete")
}

func inviteProject(c echo.Context) error {
	user := interceptor.User
	form := &InviteProjectForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(user.ID, form.ProjectId); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{constants.PermissionException})
	}
	auth, err := data.AuthByMailAddress(form.MailAddress)
	if isErr(err) {
		return CreateErrorResponse(err, c)
	}
	if err := data.UserProjectByUserIdProjectId(auth.ID, form.ProjectId); err == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"already exists user_project"})
	}
	inviteInfo := InviteInfo{ProjectId: form.ProjectId, UserId: auth.UserId}
	inviteInfoJson, err := json.Marshal(inviteInfo)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	token, _ := MakeRandomStr()
	data.RedisSet(string(inviteInfoJson), token)
	message := "http://localhost:8081/join/" + token
	if err := mail.SendMail(auth.MailAddress, message); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, SuccessResponse{"send mail to" + auth.MailAddress})
}

func joinProject(c echo.Context) error {
	form := &JoinProjectForm{}
	if err := c.Bind(form); err != nil {
		return CreateErrorResponse(err, c)
	}
	inviteInfoJson, _ := data.RedisGetInviteInfo(form.Token)

	var inviteInfo = new(InviteInfo)
	if err := json.Unmarshal([]byte(inviteInfoJson), inviteInfo); err != nil {
		return CreateErrorResponse(err, c)
	}

	if err := data.UserProjectByUserIdProjectId(inviteInfo.UserId, inviteInfo.ProjectId); err == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"already exists user_project"})
	}

	data.InsertUserProject(data.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})

	return c.JSON(http.StatusOK,
		response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})
}
