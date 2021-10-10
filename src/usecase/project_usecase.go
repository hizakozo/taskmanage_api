package usecase

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/mail"
	"taskmanage_api/src/model"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
)

type projectUsecase struct {
	pr  domain.ProjectRepository
	upr domain.UserProjectRepository
	rr  domain.RedisRepository
	ar  domain.AuthRepository
	sr  domain.StatusRepository
}

func NewProjectUsecase(pr domain.ProjectRepository, upr domain.UserProjectRepository, rr domain.RedisRepository,
	ar domain.AuthRepository, sr domain.StatusRepository) domain.ProjectUsecase {
	return &projectUsecase{
		pr:  pr,
		upr: upr,
		rr:  rr,
		ar:  ar,
		sr:  sr,
	}
}

func (pu *projectUsecase) GetList(userId int, userName, avatar string) *response.ProjectList {
	responseUser := response.User{Id: userId, Name: userName, Avatar: avatar}

	var responseProjects []response.Project
	for _, project := range pu.pr.ProjectsByUserId(userId) {
		responseProjects =
			append(responseProjects, response.Project{Id: project.ID, Name: project.ProjectName,
				Description: project.Description, Avatar: project.ProjectAvatar})
	}
	return &response.ProjectList{User: responseUser, Projects: responseProjects}
}

func (pu *projectUsecase) Create(projectName, description string, userId int) error {
	var project domain.Project
	project.ProjectName = projectName
	project.Description = description
	insertProjectId := pu.pr.InsertProject(project)

	var userProject domain.UserProject
	userProject.UserId = userId
	userProject.ProjectId = insertProjectId
	pu.upr.InsertUserProject(userProject)

	statuses := constants.Statuses
	for _, v := range statuses {
		var status domain.Status
		status.ProjectId = insertProjectId
		status.Progress = v.Progress
		status.StatusName = v.Name
		pu.sr.InsertStatus(status)
	}
	return nil
}

func (pu *projectUsecase) Update(userId, projectId int, projectName, description, avatar string, c echo.Context) error {
	if userProject := pu.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}

	project := domain.Project{ID: projectId, ProjectName: projectName,
		Description: description, ProjectAvatar: avatar}
	pu.pr.UpdateProject(project)
	return nil
}

func (pu *projectUsecase) Delete(userId, projectId int, c echo.Context) error {
	if userProject := pu.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return exception.PermissionException(c)
	}
	pu.pr.DeleteProject(projectId)
	return nil
}

func (pu *projectUsecase) Invite(userId, projectId int, mailAddress string, c echo.Context) (*response.UserProject, error) {
	if userProject := pu.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	auth, err := pu.ar.AuthByMailAddress(mailAddress)
	if err != nil {
		return nil, exception.NotFoundData(c)
	}
	if userProject := pu.upr.UserProjectByUserIdProjectId(auth.ID, projectId); len(userProject) > 0 {
		return nil, exception.DataAlreadyExists(c, "user_project")
	}
	inviteInfo := model.InviteInfo{ProjectId: projectId, UserId: auth.UserId}
	inviteInfoJson, _ := json.Marshal(inviteInfo)
	token, _ := utils.MakeRandomStr()
	pu.rr.RedisSet(string(inviteInfoJson), token)
	message :=
		constants.MailBody +
			constants.Params.FrontUrl + "#/join/" + token
	_ = mail.SendMail(auth.MailAddress, message)
	return &response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId}, nil
}

func (pu *projectUsecase) Join(token string, c echo.Context) (*response.UserProject, error) {
	inviteInfoJson, _ := pu.rr.RedisGetInviteInfo(token)

	var inviteInfo = new(model.InviteInfo)
	_ = json.Unmarshal([]byte(inviteInfoJson), inviteInfo)

	if userProject := pu.upr.UserProjectByUserIdProjectId(inviteInfo.UserId, inviteInfo.ProjectId); len(userProject) > 0 {
		return nil, exception.DataAlreadyExists(c, "user_project")
	}

	pu.upr.InsertUserProject(domain.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId})

	pu.rr.RedisDelete(token)

	return &response.UserProject{UserId: inviteInfo.UserId, ProjectId: inviteInfo.ProjectId}, nil
}
