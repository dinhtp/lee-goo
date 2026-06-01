package validate

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// BindAndValidate binds the request body/query/params into v and returns
// a 400 HTTPError on bind failure. Validation tags are intentionally left
// to the caller so each module can apply its own validator.
func BindAndValidate(c echo.Context, v any) error {
	if err := c.Bind(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
