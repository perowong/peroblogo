package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/perowong/peroblogo/conf"
)

const (
	gitAccessTokenUrl = "https://github.com/login/oauth/access_token"
	gitUserUrl        = "https://api.github.com/user"
)

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUser struct {
	OpenID    int    `json:"id"`
	Nickname  string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
}

func GetGithubAccessToken(code string, result *GithubAccessToken) error {
	param := make(map[string]string)
	param["code"] = code
	param["client_id"] = conf.C.Github.ClientID
	param["client_secret"] = conf.C.Github.ClientSecret
	paramBytes, _ := json.Marshal(param)
	body := strings.NewReader(string(paramBytes))

	req, _ := http.NewRequest("POST", gitAccessTokenUrl, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}

func GetGithubUser(githubToken GithubAccessToken, result *GithubUser) error {
	req, _ := http.NewRequest("GET", gitUserUrl, nil)
	req.Header.Add("Authorization", githubToken.TokenType+" "+githubToken.AccessToken)
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}
