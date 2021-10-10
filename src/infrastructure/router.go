package infrastructure

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"taskmanage_api/src/controller"
	"taskmanage_api/src/interceptor"
	"taskmanage_api/src/repository"
	"taskmanage_api/src/usecase"
)

var Echo *echo.Echo

func init() {
	e := echo.New()
	e.Use(middleware.CORS())

	jwtMiddleWare := middleware.JWT([]byte("secret"))

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
	e.POST("/signIn", func(c echo.Context) error { return uc.SignIn(c) })
	e.POST("/signUp", func(c echo.Context) error { return uc.SignUp(c) })
	u := e.Group("/user")
	u.Use(jwtMiddleWare)
	u.GET("/:project_id", func(c echo.Context) error { return uc.GetUsersInProject(c) }, i.CsrfAuth)
	u.POST("/signOut", func(c echo.Context) error { return uc.SignOut(c) })
	u.GET("/profile", func(c echo.Context) error { return uc.GetUserProfile(c) }, i.CsrfAuth)

	tu := usecase.NewTicketUsecase(ur, tr, tsr, pr, upr, sr)
	tc := controller.NewTicketController(tu)
	t := e.Group("/tickets")
	t.Use(jwtMiddleWare)
	t.POST("/create", func(c echo.Context) error { return tc.Create(c) }, i.CsrfAuth)
	t.GET("/:project_id", func(c echo.Context) error { return tc.GetList(c) }, i.CsrfAuth)
	t.PUT("/update", func(c echo.Context) error { return tc.Update(c) }, i.CsrfAuth)
	t.PUT("/status", func(c echo.Context) error { return tc.ChangeStatus(c) }, i.CsrfAuth)
	t.GET("/detail/:ticket_id", func(c echo.Context) error { return tc.Detail(c) }, i.CsrfAuth)
	t.DELETE("/delete/:ticket_id", func(c echo.Context) error { return tc.Delete(c) }, i.CsrfAuth)

	cu := usecase.NewCommentUsecase(ur, tr, cr)
	cc := controller.NewCommentController(cu)
	c := e.Group("/comments")
	c.Use(jwtMiddleWare)
	c.GET("/:ticket_id", func(c echo.Context) error { return cc.Detail(c) }, i.CsrfAuth)
	c.POST("", func(c echo.Context) error { return cc.Create(c) }, i.CsrfAuth)

	pu := usecase.NewProjectUsecase(pr, upr, rr, ar, sr)
	pc := controller.NewProjectController(pu)
	p := e.Group("/projects")
	p.Use(jwtMiddleWare)
	p.GET("", func(c echo.Context) error { return pc.GetList(c) }, i.CsrfAuth)
	p.POST("/create", func(c echo.Context) error { return pc.Create(c) }, i.CsrfAuth)
	p.PUT("/update", func(c echo.Context) error { return pc.Update(c) }, i.CsrfAuth)
	p.DELETE("/:project_id", func(c echo.Context) error { return pc.Delete(c) }, i.CsrfAuth)
	p.POST("/invite", func(c echo.Context) error { return pc.Invite(c) }, i.CsrfAuth)
	p.POST("/join", func(c echo.Context) error { return pc.Join(c) }, i.CsrfAuth)

	su := usecase.NewStatusUsecase(tr, upr, sr)
	sc := controller.NewStatusController(su)
	s := e.Group("/statuses")
	s.Use(jwtMiddleWare)
	s.GET("/:project_id", func(c echo.Context) error { return sc.GetList(c) }, i.CsrfAuth)
	s.POST("/create", func(c echo.Context) error { return sc.Create(c) }, i.CsrfAuth)
	s.PUT("/update", func(c echo.Context) error { return sc.Update(c) }, i.CsrfAuth)
	s.DELETE("/delete/:status_id", func(c echo.Context) error { return sc.Delete(c) }, i.CsrfAuth)

	Echo = e
}

func AppRun() {
	Echo.Logger.Fatal(Echo.Start(":1313"))
}
