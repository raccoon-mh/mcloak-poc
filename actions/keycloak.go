package actions

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"

	"github.com/Nerzal/gocloak/v13"
)

// HomeHandler is a default handler to serve up
// a home page.
func KcHomeHandler(c buffalo.Context) error {
	c.Set("simplestr", "welcome mcloak home!")
	return c.Render(http.StatusOK, r.HTML("kctest/index.html"))
}

func KcCreateUserHandler(c buffalo.Context) error {
	var KC_admin = os.Getenv("KC_admin")
	var KC_passwd = os.Getenv("KC_passwd")
	var KC_uri = os.Getenv("KC_uri")
	// fmt.Println("KC ENV :", KC_admin, KC_passwd, KC_uri)

	client := gocloak.NewClient(KC_uri)
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, KC_admin, KC_passwd, "master")
	if err != nil {
		c.Set("simplestr", err.Error()+"### Something wrong with the credentials or url ###")
		return c.Render(http.StatusOK, r.HTML("kctest/index.html"))
		// panic("Something wrong with the credentials or url")
	}

	fmt.Println(token)

	user := gocloak.User{
		FirstName: gocloak.StringP("ra"),
		LastName:  gocloak.StringP("ccoon"),
		Email:     gocloak.StringP("mega@zone.cloud"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("raccoon"),
	}

	_, err = client.CreateUser(ctx, token.AccessToken, "master", user)
	if err != nil {
		c.Set("simplestr", err.Error()+"### Oh no!, failed to create user :( ###")
		return c.Render(http.StatusOK, r.HTML("kctest/index.html"))
		// panic("Oh no!, failed to create user :(")
	}

	c.Set("simplestr", "success")
	return c.Render(http.StatusOK, r.HTML("kctest/index.html"))
}

func KcLoginHandler(c buffalo.Context) error {
	var KC_admin = os.Getenv("KC_admin")
	var KC_passwd = os.Getenv("KC_passwd")
	var KC_uri = os.Getenv("KC_uri")
	// fmt.Println("KC ENV :", KC_admin, KC_passwd, KC_uri)

	client := gocloak.NewClient(KC_uri)
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, KC_admin, KC_passwd, "master")
	if err != nil {
		c.Set("simplestr", err.Error()+"### Something wrong with the credentials or url ###")
		return c.Render(http.StatusOK, r.HTML("kctest/index.html"))
		// panic("Something wrong with the credentials or url")
	}

	return c.Render(http.StatusOK, r.JSON(token))
}
