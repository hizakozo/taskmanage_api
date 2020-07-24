package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"taskmanage_api/src/handler"
	"taskmanage_api/src/interceptor"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/user/signIn", handler.SignIn)
	e.POST("/user/signUp", handler.SignUp)
	e.GET("/user/:project_id", handler.GetUsersInProject, interceptor.CsrfAuth)
	e.POST("/user/signOut", handler.SignOut)
	e.GET("/user/profile", handler.GetUserProfile, interceptor.CsrfAuth)

	e.GET("/projects", handler.GetProjectList, interceptor.CsrfAuth)
	e.POST("/projects", handler.CreateProject, interceptor.CsrfAuth)
	e.PUT("/projects", handler.UpdateProject, interceptor.CsrfAuth)
	e.DELETE("/projects/:project_id", handler.DeleteProject, interceptor.CsrfAuth)
	e.POST("/projects/invite", handler.InviteProject, interceptor.CsrfAuth)
	e.POST("/projects/join", handler.JoinProject)

	e.POST("/tickets", handler.CreateTicket, interceptor.CsrfAuth)
	e.GET("/tickets/:project_id", handler.GetTicketList, interceptor.CsrfAuth)
	e.PUT("/tickets", handler.UpdateTicket, interceptor.CsrfAuth)
	e.PUT("/tickets/status", handler.ChangeStatus, interceptor.CsrfAuth)
	e.GET("/tickets/detail/:ticket_id", handler.GetTicketDetail, interceptor.CsrfAuth)
	e.DELETE("/tickets/delete/:ticket_id", handler.DeleteTicket, interceptor.CsrfAuth)

	e.GET("/statuses/:project_id", handler.GetStatusList, interceptor.CsrfAuth)
	e.POST("/statuses", handler.CreateStatus, interceptor.CsrfAuth)
	e.PUT("/statuses", handler.UpdateStatus, interceptor.CsrfAuth)
	e.DELETE("/statuses/delete/:status_id", handler.DeleteStatus, interceptor.CsrfAuth)

	e.GET("/comments/:ticket_id", handler.GetComment, interceptor.CsrfAuth)
	e.POST("/comments", handler.CreateComment, interceptor.CsrfAuth)

	e.Logger.Fatal(e.Start(":1313"))
}