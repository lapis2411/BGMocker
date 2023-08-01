package receive

import (
	"io"

	"github.com/labstack/echo/v4"
)

func ReceiveFile(c echo.Context, name string) ([]byte, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return io.ReadAll(src)
}
