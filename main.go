package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := echo.New()

	e.POST("/signIn", login)
	e.POST("/signUp", signUp)

	e.GET("/projects", getProjectList)
	e.POST("/projects", createProject)
	e.PUT("/projects", updateProject)
	e.DELETE("/projects/:project_id", deleteProject)
	e.POST("/projects/invite", inviteProject)
	e.POST("/projects/join", joinProject)


	e.POST("/tickets", createTicket)
	e.GET("/tickets/:project_id", getTicketList)
	e.PUT("/tickets", updateTicket)
	e.PUT("/tickets/status", changeStatus)
	e.GET("/tickets/detail/:ticket_id", displayTicketDetail)
	e.DELETE("/tickets/:ticket_id", deleteTicket)

	e.GET("/statuses/:project_id", getStatusList)
	e.POST("/statuses", createStatus)
	e.PUT("/statuses", updateStatus)
	e.DELETE("/statuses/delete/:status_id", deleteStatus)

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
