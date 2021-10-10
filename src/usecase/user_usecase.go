package usecase

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"taskmanage_api/src/constants"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
	"taskmanage_api/src/response"
	"taskmanage_api/src/utils"
	"time"
)

type userUseCase struct {
	ur  domain.UserRepository
	ar  domain.AuthRepository
	rr  domain.RedisRepository
	upr domain.UserProjectRepository
	pr  domain.ProjectRepository
	tr  domain.TicketRepository
}

func NewUserUseCase(ur domain.UserRepository, ar domain.AuthRepository,
	rr domain.RedisRepository, upr domain.UserProjectRepository,
	pr domain.ProjectRepository, tr domain.TicketRepository) domain.UserUseCase {
	return &userUseCase{
		ur:  ur,
		ar:  ar,
		rr:  rr,
		upr: upr,
		pr:  pr,
		tr:  tr,
	}
}

func (uu *userUseCase) SignIn(loginId, password string) (*response.LoginResponse, error) {
	auth, err := uu.ar.AuthByLoginId(loginId)
	if err != nil {
		return nil, err
	}
	if err = utils.PasswordVerify(auth.Password, password); err != nil {
		return nil, err
	}
	user, err := uu.ur.UserById(auth.UserId)
	if err != nil {
		return nil, err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Minute * 24).Unix()
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	return &response.LoginResponse{UserToken: t, UserId: user.ID}, nil
}

func (uu *userUseCase) SignUp(name, loginId, password, mailAddress string, avatar *multipart.FileHeader, c echo.Context) error {
	if _, err := uu.ar.AuthByLoginId(loginId); err == nil {
		return exception.DataAlreadyExists(c, "login id")
	}
	if _, err := uu.ar.AuthByMailAddress(mailAddress); err == nil {
		return exception.DataAlreadyExists(c, "mail address")
	}
	user := domain.User{Name: name}
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
	id := uu.ur.InsertUser(user)
	safetyPass := utils.CreateSafetyPass(password)
	auth := domain.Auth{UserId: id, LoginId: loginId, Password: safetyPass, MailAddress: mailAddress}

	uu.ar.InsertAuth(auth)

	return nil
}

func (uu *userUseCase) GetUsersInProject(userId, projectId int, c echo.Context) (*response.UserList, error) {
	if userProject := uu.upr.UserProjectByUserIdProjectId(userId, projectId); len(userProject) == 0 {
		return nil, exception.PermissionException(c)
	}
	users := uu.ur.UserByProjectId(projectId)
	var responseUsers []response.User
	for _, user := range users {
		responseUser := response.User{Id: user.ID, Name: user.Name}
		if user.Avatar != "" {
			responseUser.Avatar = constants.Params.S3Url + user.Avatar
		}
		responseUsers = append(responseUsers, responseUser)
	}
	return &response.UserList{Users: responseUsers}, nil
}

func (uu *userUseCase) SignOut(token string, userId int, c echo.Context) error {
	user, _ := uu.rr.RedisGet(token)
	if user.ID != userId {
		return exception.PermissionException(c)
	}
	uu.rr.RedisDelete(token)
	return nil
}

func (uu *userUseCase) GetUserProfile(userId int, c echo.Context) (*response.UserProfile, error) {
	user, err := uu.ur.UserById(userId)
	auth, err := uu.ar.AuthByUserId(userId)
	if err != nil {
		return nil, exception.NotFoundData(c)
	}
	projects := uu.pr.ProjectsByUserId(userId)
	responseUser := response.User{Id: user.ID, Name: user.Name,
		Avatar: user.Avatar, MailAddress: auth.MailAddress}

	var responseProjects []response.ProjectTickets
	for _, project := range projects {
		var responseTickets []response.Ticket
		for _, ticket := range uu.tr.TicketByProjectIdWorker(project.ID, user.ID) {
			responseTickets = append(responseTickets, response.Ticket{Id: ticket.ID, Title: ticket.Title})
		}
		responseProjects = append(responseProjects, response.ProjectTickets{Id: project.ID, Name: project.ProjectName,
			Tickets: responseTickets})
	}

	return &response.UserProfile{User: responseUser, Projects: responseProjects}, nil
}
