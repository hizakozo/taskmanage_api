package main

import (
	"./interceptor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/user/signIn", login)
	e.POST("/user/signUp", signUp)
	e.GET("/user/:project_id", getUsersInProject, interceptor.CsrfAuth)

	e.GET("/projects", getProjectList, interceptor.CsrfAuth)
	e.POST("/projects", createProject, interceptor.CsrfAuth)
	e.PUT("/projects", updateProject, interceptor.CsrfAuth)
	e.DELETE("/projects/:project_id", deleteProject, interceptor.CsrfAuth)
	e.POST("/projects/invite", inviteProject, interceptor.CsrfAuth)
	e.POST("/projects/join", joinProject, interceptor.CsrfAuth)

	e.POST("/tickets", createTicket, interceptor.CsrfAuth)
	e.GET("/tickets/:project_id", getTicketList, interceptor.CsrfAuth)
	e.PUT("/tickets", updateTicket, interceptor.CsrfAuth)
	e.PUT("/tickets/status", changeStatus, interceptor.CsrfAuth)
	e.GET("/tickets/detail/:ticket_id", displayTicketDetail, interceptor.CsrfAuth)
	e.DELETE("/tickets/delete/:ticket_id", deleteTicket, interceptor.CsrfAuth)

	e.GET("/statuses/:project_id", getStatusList, interceptor.CsrfAuth)
	e.POST("/statuses", createStatus, interceptor.CsrfAuth)
	e.PUT("/statuses", updateStatus, interceptor.CsrfAuth)
	e.DELETE("/statuses/delete/:status_id", deleteStatus, interceptor.CsrfAuth)

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})

	e.Logger.Fatal(e.Start(":1313"))
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func CreateErrorResponse(err error, c echo.Context) error {
	log.Error(err)
	return c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
}

func isErr(err error) bool {
	return err != nil
}