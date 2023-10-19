package actions

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/versent/saml2aws/v2/pkg/cfg"
	"github.com/versent/saml2aws/v2/pkg/creds"
)

// HomeHandler is a default handler to serve up
// a home page.
func LoginHandler(c buffalo.Context) error {
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, r.HTML("login/index.html"))
	}

	var username = c.Request().FormValue("username")
	var password = c.Request().FormValue("password")

	token, err := KC_client.Login(c, KC_clientID, KC_clientSecret, KC_realm, username, password)
	if err != nil {
		c.Set("simplestr", err.Error())
		return c.Render(http.StatusOK, r.HTML("simplestr.html"))
	}
	userinfo, err := KC_client.GetUserInfo(c, token.AccessToken, KC_realm)
	if err != nil {
		c.Set("simplestr", err.Error())
		return c.Render(http.StatusOK, r.HTML("simplestr.html"))
	}
	fmt.Println("userinfo", userinfo)

	c.Session().Set("AccessToken", token.AccessToken)
	// c.Session().Set("RefreshToken", token.RefreshToken) // TODO : save in db to reduce cookie

	return c.Redirect(302, "/buffalo/authuser")
}

func AuthUserTestPageHandler(c buffalo.Context) error {
	c.Set("simplestr", "You are good to go")
	return c.Render(http.StatusOK, r.HTML("simplestr.html"))
}

func NotAuthUserTestPageHandler(c buffalo.Context) error {
	c.Set("simplestr", "You are blocked by middleware")
	return c.Render(http.StatusOK, r.HTML("simplestr.html"))
}

func Testsaml(c buffalo.Context) error {

	account := &cfg.IDPAccount{
		URL:                  os.Getenv("SAML_IDP_Initiated_URL"),
		Username:             os.Getenv("SAML_user"),
		Provider:             "KeyCloak",
		MFA:                  "Auto",
		SkipVerify:           false,
		AmazonWebservicesURN: "urn:amazon:webservices",
		SessionDuration:      900,
		Profile:              "saml",
		RoleARN:              "",
		Region:               "",
	}

	loginDetails := &creds.LoginDetails{
		URL:          account.URL,
		Username:     account.Username,
		Password:     os.Getenv("SAML_password"),
		MFAToken:     "",
		DuoMFAOption: "",
	}

	cred, err := Login(account, loginDetails)
	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	fmt.Println("########## AWSAccessKey  : ", cred.AWSAccessKey)
	fmt.Println("########## AWSSecretKey  : ", cred.AWSSecretKey)
	fmt.Println("########## AWSSecurityToken  : ", cred.AWSSecurityToken)
	fmt.Println("########## AWSSessionToken  : ", cred.AWSSessionToken)
	fmt.Println("########## Expires  : ", cred.Expires.Format(time.RFC3339))

	return c.Render(http.StatusOK, r.JSON("success"))
}
