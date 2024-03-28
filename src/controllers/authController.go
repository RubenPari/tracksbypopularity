package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/RubenPari/tracksByPopularity/src/utils"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Login(c echo.Context) error {
	state := utils.RandomString(16)

	// save state in session
	sess, _ := session.Get("session", c)
	sess.Values["state"] = state
	_ = sess.Save(c.Request(), c.Response())

	redirectURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		os.Getenv("CLIENT_ID"),
		os.Getenv("SCOPES"),
		os.Getenv("REDIRECT_URI"),
		state)

	return c.Redirect(302, redirectURL)
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

func Redirect(c echo.Context) error {
	// get code and state from query params
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	// initialize session
	sess, _ := session.Get("session", c)

	// check state
	if sess.Values["state"] != state {
		return c.String(400, "Invalid state")
	}

	// prepare request body for access token
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", os.Getenv("REDIRECT_URI"))

	// create request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID")+":"+os.Getenv("CLIENT_SECRET"))))

	// send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// close response body
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// read response bytes
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// parse access token
	accessToken := &AccessToken{}
	err = json.Unmarshal(bodyBytes, accessToken)
	if err != nil {
		return err
	}

	// save access token in session
	sess.Values["access_token"] = accessToken.AccessToken

	// save session
	_ = sess.Save(c.Request(), c.Response())

	return c.String(200, "Logged in successfully")
}

func Logout(c echo.Context) error {
	// initialize session
	sess, _ := session.Get("session", c)

	// delete access token
	delete(sess.Values, "access_token")

	// save session
	_ = sess.Save(c.Request(), c.Response())

	return c.String(200, "Logged out successfully")
}
