package app

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lapis2411/board-game-maker/app/image"
	"github.com/lapis2411/board-game-maker/app/receive"
)

func main() {
	e := echo.New()
	e.GET("/users/:id", getUser)
	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t
	e.GET("/top", func(c echo.Context) error {
		return c.Render(http.StatusOK, "top.html", nil)
	})
	e.POST("/generate", generate)
	e.GET("/*", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/top")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func generate(c echo.Context) error {
	style, err := receive.ReceiveFile(c, "styleCSV")
	if err != nil {
		return fmt.Errorf("failed to receive styleCSV: %w", err)
	}
	card, err := receive.ReceiveFile(c, "cardCSV")
	if err != nil {
		return fmt.Errorf("failed to receive cardCSV: %w", err)
	}
	res, err := image.PrintsJsons(style, card)
	// base64形式の画像のリストをレスポンスとして返す
	return c.JSON(http.StatusOK, res)
}
