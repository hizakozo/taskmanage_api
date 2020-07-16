package interceptor

import (
	"taskmanage_api/constants"
	"taskmanage_api/data"
	"github.com/labstack/echo"
	"net/http"
)

var User = data.User{}

func CsrfAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := data.RedisGet(c.Request().Header.Get("user_token"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, constants.TokenFailed)
		}
		User = user
		if err := next(c); err != nil {
			c.Error(err)
		}
		User = data.User{}
		return nil
	}
}