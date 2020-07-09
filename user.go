package main

import (
	"./data"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type TokenResponse struct {
	UserToken string `json:"user_token"`
}

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
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	fmt.Println(auth)

	if err := passwordVerify(auth.Password, form.Password); err != nil {
		return CreateErrorResponse(err, c)
	}

	user, err := data.UserById(auth.UserId)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	var key, _ = MakeRandomStr()

	userJson, _ := json.Marshal(user)

	data.RedisSet(string(userJson), key)

	response := TokenResponse{UserToken: key}

	return c.JSON(http.StatusOK, response)
}
