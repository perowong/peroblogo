package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/api"
	"github.com/perowong/peroblogo/model"
)

type LoginReq struct {
	Code string `binding:"required"`
	Type string `binding:"required"`
}

type LoginResp struct {
	*model.User
	Token string
}

// user login
func Login(c *gin.Context) {
	var req LoginReq
	var err error

	ctxUtils := &GinCtxUtils{Context: c}
	if !ctxUtils.GetReqObject(&req) {
		return
	}

	// req.Type == "github"
	var githubAccessTokenResp api.GithubAccessToken
	err = api.GetGithubAccessToken(req.Code, &githubAccessTokenResp)

	if err != nil {
		log.Println("getGithubAccessToken error: ", err.Error())
		ctxUtils.ReplyFailParam()
		return
	}

	var githubUserResp api.GithubUser
	err = api.GetGithubUser(githubAccessTokenResp, &githubUserResp)
	if err != nil {
		log.Println("getGithubUser error: ", err.Error())
		ctxUtils.ReplyFailParam()
		return
	}

	if githubUserResp.OpenID == 0 {
		log.Println("githubUserResp error: ", githubUserResp)
		ctxUtils.ReplyFail(ErrCodeParam, "bad code")
		return
	}

	log.Printf("githubUserResp: %#v", githubUserResp)
	daoObj := model.NewModel()
	userID, err := daoObj.GetUserIDBy(strconv.Itoa(githubUserResp.OpenID))
	if err != nil {
		log.Println(err.Error())
		ctxUtils.ReplyFailServer()
		return
	}

	if userID == 0 {
		userID, err = daoObj.AddUser(&model.User{
			OpenID:    strconv.Itoa(githubUserResp.OpenID),
			AuthType:  "github",
			Nickname:  githubUserResp.Nickname,
			AvatarUrl: githubUserResp.AvatarUrl,
			Email:     githubUserResp.Email,
		})

		if err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailServer()
			return
		}
	}

	var token string
	exist, err := daoObj.CheckUserTokenExistBy(userID)
	if err != nil {
		log.Println(err.Error())
		ctxUtils.ReplyFailServer()
		return
	}
	if !exist {
		token, err = daoObj.AddUserToken(userID)
		if err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailServer()
			return
		}
	} else {
		token, err = daoObj.UpdateUserTokenByUserID(userID)
		if err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailServer()
			return
		}
	}

	ctxUtils.ReplyOk(&LoginResp{
		User: &model.User{
			ID:        userID,
			OpenID:    strconv.Itoa(githubUserResp.OpenID),
			AuthType:  "github",
			Nickname:  githubUserResp.Nickname,
			AvatarUrl: githubUserResp.AvatarUrl,
			Email:     githubUserResp.Email,
		},
		Token: token,
	})
}
