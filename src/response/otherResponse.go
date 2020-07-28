package response

import (
	"github.com/labstack/echo"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func CreateErrorResponse(err error, c echo.Context) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
}
