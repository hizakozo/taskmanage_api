package infrastructure

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"taskmanage_api/src/controller"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/repository"
	"taskmanage_api/src/usecase"
)

var Echo *echo.Echo

func init() {
	e := echo.New()
	e.Use(middleware.CORS())

	rr := repository.NewRedisRepository(Redis)
	ur := repository.NewUserRepository(Db)
	tr := repository.NewTicketRepository(Db)
	tsr := repository.NewTicketStatusRepository(Db)
	pr := repository.NewProjectRepository(Db)
	upr := repository.NewUserProjectRepository(Db)
	sr := repository.NewStatusRepository(Db)
	ar := repository.NewAuthRepository(Db)
	cr := repository.NewCommentRepository(Db)

	i := interceptor.NewIntercept(rr)

	uu := usecase.NewUserUseCase(ur, ar, rr, upr, pr, tr)
	uc := controller.NewUserController(uu)
	e.POST("/user/signIn", func(c echo.Context) error { return uc.SignIn(c) })
	e.POST("/user/signUp", func(c echo.Context) error { return uc.SignUp(c) })
	e.GET("/user/:project_id", func(c echo.Context) error { return uc.GetUsersInProject(c) }, i.CsrfAuth)
	e.POST("/user/signOut", func(c echo.Context) error { return uc.SignOut(c) })
	e.GET("/user/profile", func(c echo.Context) error { return uc.GetUserProfile(c) }, i.CsrfAuth)

	tu := usecase.NewTicketUsecase(ur, tr, tsr, pr, upr, sr)
	tc := controller.NewTicketController(tu)
	e.POST("/tickets/create", func(c echo.Context) error { return tc.Create(c) }, i.CsrfAuth)
	e.GET("/tickets/:project_id", func(c echo.Context) error { return tc.GetList(c) }, i.CsrfAuth)
	e.PUT("/tickets/update", func(c echo.Context) error { return tc.Update(c) }, i.CsrfAuth)
	e.PUT("/tickets/status", func(c echo.Context) error { return tc.ChangeStatus(c) }, i.CsrfAuth)
	e.GET("/tickets/detail/:ticket_id", func(c echo.Context) error { return tc.Delete(c) }, i.CsrfAuth)
	e.DELETE("/tickets/delete/:ticket_id", func(c echo.Context) error { return tc.Delete(c) }, i.CsrfAuth)

	cu := usecase.NewCommentUsecase(ur, tr, cr)
	cc := controller.NewCommentController(cu)
	e.GET("/comments/:ticket_id", func(c echo.Context) error { return cc.Detail(c) }, i.CsrfAuth)
	e.POST("/comments", func(c echo.Context) error { return cc.Create(c) }, i.CsrfAuth)

	pu := usecase.NewProjectUsecase(pr, upr, rr, ar, sr)
	pc := controller.NewProjectController(pu)
	e.GET("/projects", func(c echo.Context) error { return pc.GetList(c) }, i.CsrfAuth)
	e.POST("/projects/create", func(c echo.Context) error { return pc.Create(c) }, i.CsrfAuth)
	e.PUT("/projects/update", func(c echo.Context) error { return pc.Update(c) }, i.CsrfAuth)
	e.DELETE("/projects/:project_id", func(c echo.Context) error { return pc.Delete(c) }, i.CsrfAuth)
	e.POST("/projects/invite", func(c echo.Context) error { return pc.Invite(c) }, i.CsrfAuth)
	e.POST("/projects/join", func(c echo.Context) error { return pc.Join(c) }, i.CsrfAuth)

	su := usecase.NewStatusUsecase(tr, upr, sr)
	sc := controller.NewStatusController(su)
	e.GET("/statuses/:project_id", func(c echo.Context) error { return sc.GetList(c) }, i.CsrfAuth)
	e.POST("/statuses/create", func(c echo.Context) error { return sc.Create(c) }, i.CsrfAuth)
	e.PUT("/statuses/update", func(c echo.Context) error { return sc.Update(c) }, i.CsrfAuth)
	e.DELETE("/statuses/delete/:status_id", func(c echo.Context) error { return sc.Delete(c) }, i.CsrfAuth)

	Echo = e
}

func Run() {
	Echo.Logger.Fatal(Echo.Start(":1313"))
}
