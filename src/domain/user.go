package domain

import (
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"taskmanage_api/src/response"
)

type User struct {
	ID       int    `gorm:"column:user_id;PRIMARY_KEY"`
	Name     string `gorm:"column:user_name"`
	Avatar   string
	Isdelete int `gorm:"column:is_delete;default:'galeone'"`
}

type SignUpForm struct {
	Name        string `form:"name" validate:"required"`
	LoginId     string `form:"login_id"  validate:"required"`
	Password    string `form:"password" validate:"required"`
	MailAddress string `form:"mail_address" validate:"required,email"`
	Avatar      string `form:"avatar"`
}

type LoginForm struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LogOutForm struct {
	UserId int `json:"user_id"`
}

type UserUseCase interface {
	SignIn(loginId, password string) (*response.LoginResponse, error)
	SignUp(name, loginId, password, mailAddress string, avatar *multipart.FileHeader, c echo.Context) error
	GetUsersInProject(userId, projectId int, c echo.Context) (*response.UserList, error)
	SignOut(token string, userId int, c echo.Context) error
	GetUserProfile(userId int, c echo.Context) (*response.UserProfile, error)
}

type UserRepository interface {
	UserById(userId int) (User, error)
	InsertUser(user User) int
	UserByProjectId(projectId int) []User
}