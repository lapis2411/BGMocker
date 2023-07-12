package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
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
	e.POST("/render", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/top")
	})
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

func getConvertFile(c echo.Context) error {
	file, err := c.FormFile("mockFile")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	return c.Redirect(http.StatusMovedPermanently, "/someurl")
}
