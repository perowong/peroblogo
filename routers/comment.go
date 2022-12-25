package routers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/dao"
	"github.com/perowong/peroblogo/utils"
)

type AddCommentReq struct {
	BlogID       string `binding:"required"`
	ParentID     int64
	ReplyID      int64
	FromUid      string `binding:"required"`
	FromNickname string `binding:"required"`
	FromEmail    string `binding:"required,email"`
	ToUid        string
	ToNickname   string
	ToEmail      string
	Content      string `binding:"required,min=1,max=600"`
}

type AddCommentResp struct {
	ID int64
	Ct time.Time
}

// Add comment
func AddComment(c *gin.Context) {
	var req AddCommentReq
	var err error

	ctxUtils := &utils.GinCtxUtils{Context: c}
	if !ctxUtils.GetReqObject(&req) {
		return
	}

	daoComment := &dao.Comment{
		BlogID:       req.BlogID,
		ParentID:     req.ParentID,
		ReplyID:      req.ReplyID,
		FromUid:      req.FromUid,
		FromNickname: req.FromNickname,
		FromEmail:    req.FromEmail,
		ToUid:        req.ToUid,
		ToNickname:   req.ToNickname,
		ToEmail:      req.ToEmail,
		Content:      req.Content,
	}

	daoObj := dao.NewDao()

	if daoComment.ReplyID != 0 {
		cResp, err := daoObj.CheckExistByID(daoComment.ReplyID)
		if cResp.ID == 0 || err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailParam()
			return
		}
	}

	var parentComment *dao.Comment
	if daoComment.ParentID != 0 {
		parentComment, err = daoObj.ReadComment(daoComment.ParentID)
		if parentComment.ID == 0 || err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailParam()
			return
		}
	}

	id, err := daoObj.AddComment(daoComment)
	if err != nil {
		log.Println(err.Error())
		ctxUtils.ReplyFailServer()
		return
	}

	if parentComment != nil {
		err = daoObj.UpdateSubCount(parentComment.ID, parentComment.SubCount+1)
		if err != nil {
			log.Println(err.Error())
			ctxUtils.ReplyFailServer()
			return
		}
	}

	ctxUtils.ReplyOk(&AddCommentResp{
		ID: id,
		Ct: time.Now(),
	})
}

type ListCommentReq struct {
	BlogID string `binding:"required"`
}

type ListCommentResp struct {
	List []*dao.Comment
}

// Query comment list
func ListComment(c *gin.Context) {
	var req ListCommentReq
	ctxUtils := &utils.GinCtxUtils{Context: c}
	if !ctxUtils.GetReqObject(&req) {
		return
	}

	daoObj := dao.NewDao()

	list, err := daoObj.ListCommentByBlogID(req.BlogID)
	if err != nil {
		log.Println(err.Error())
		ctxUtils.ReplyFailServer()
		return
	}

	if len(list) == 0 {
		ctxUtils.ReplyOk(&ListCommentResp{List: make([]*dao.Comment, 0)})
		return
	}

	for _, item := range list {
		if item.SubCount != 0 {
			subList, err := daoObj.ListCommentByParentID(item.ID)
			if err != nil {
				log.Println(err.Error())
				ctxUtils.ReplyFailServer()
				return
			}
			item.Children = subList
		}
	}

	ctxUtils.ReplyOk(&ListCommentResp{List: list})
}
