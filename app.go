package main

import (
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("static/template/*.html")),
	}

	e.Renderer = t

	// set up session
	store := session.NewCookieStore([]byte("secret-key"))
	store.MaxAge(86400) // log out after 1 day
	e.Use(session.Sessions("ESESSION", store))

	e.GET("/login", ShowLoginHtml)
	e.POST("/login", Login)

	e.Logger.Fatal(e.Start(":9000"))
}

type LoginForm struct {
	UserId       string
	Password     string
	ErrorMessage string
}

type CompleteJson struct {
	Success bool `json:"success"`
}

func ShowLoginHtml(c echo.Context) error {
	session := session.Default(c)

	loginId := session.Get("loginCompleted")
	if loginId != nil && loginId == "completed" {
		completeJson := CompleteJson{
			Success: true,
		}
		return c.JSON(http.StatusOK, completeJson)
	}
	fmt.Println("rendering..")
	return c.Render(http.StatusOK, "login", LoginForm{})
}

func Login(c echo.Context) error {
	loginForm := LoginForm{
		UserId:   c.FormValue("userId"),
		Password: c.FormValue("password"),
	}

	userId := html.EscapeString(loginForm.UserId)
	password := html.EscapeString(loginForm.Password)
	fmt.Println(loginForm.UserId, loginForm.Password)
	fmt.Println(userId, password)

	if userId != "userId" && password != "password" {
		loginForm.ErrorMessage = "Oops!"
		return c.Render(http.StatusOK, "login", loginForm)
	}

	// save login success to session
	session := session.Default(c)
	session.Set("loginCompleted", "completed")
	session.Save()

	completeJson := CompleteJson{
		Success: true,
	}

	return c.JSON(http.StatusOK, completeJson)
}
