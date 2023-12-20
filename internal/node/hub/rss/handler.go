package rss

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-node/common/http/response"
	"github.com/naturalselectionlabs/rss3-node/schema"
)

type Response struct {
	Data []*schema.Feed `json:"data"`
}

// GetRSSHubHandler get rsshub data from rsshub node
func (h *Hub) GetRSSHubHandler(c echo.Context) error {
	path := c.Param("*")
	rawQuery := c.Request().URL.RawQuery

	data, err := h.getRSSHubData(c.Request().Context(), path, rawQuery)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSON(http.StatusOK, Response{
		Data: data,
	})
}
