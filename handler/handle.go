package handle

import (
	"github.com/labstack/echo"
	"net/http"
)

/* func */
func Index(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello World")
}
