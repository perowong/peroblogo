package controller

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/model"
)

type AddCommentReq struct {
	BlogID       string `binding:"required"`
	ParentID     int64
	ReplyID      int64
	FromUid      int64  `binding:"required"`
	FromNickname string `binding:"required"`
	FromAvatar   string
	ToUid        int64
	ToNickname   string
	ToAvatar     string
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

	ctxUtils := &GinCtxUtils{Context: c}
	if !ctxUtils.GetReqObject(&req) {
		return
	}

	daoComment := &model.Comment{
		BlogID:       req.BlogID,
		ParentID:     req.ParentID,
		ReplyID:      req.ReplyID,
		FromUid:      req.FromUid,
		FromNickname: req.FromNickname,
		FromAvatar:   req.FromAvatar,
		ToUid:        req.ToUid,
		ToNickname:   req.ToNickname,
		ToAvatar:     req.ToAvatar,
		Content:      req.Content,
	}

	daoObj := model.NewModel()

	if daoComment.ReplyID != 0 {
		exist, err := daoObj.CheckCommentExistBy(daoComment.ReplyID)
		if !exist || err != nil {
			log.Printf("%#v", err)
			ctxUtils.ReplyFailParam()
			return
		}
	}

	var parentComment *model.Comment
	if daoComment.ParentID != 0 {
		parentComment, err = daoObj.ReadComment(daoComment.ParentID)
		if parentComment.ID == 0 || parentComment.ParentID != 0 || err != nil {
			log.Printf("%#v", err)
			ctxUtils.ReplyFailParam()
			return
		}
	}

	id, err := daoObj.AddComment(daoComment)
	if err != nil {
		log.Printf("%#v", err)
		ctxUtils.ReplyFailServer()
		return
	}

	if parentComment != nil {
		err = daoObj.UpdateSubCount(parentComment.ID, parentComment.SubCount+1)
		if err != nil {
			log.Printf("%#v", err)
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
	List []*model.Comment
}

// Query comment list
func ListComment(c *gin.Context) {
	var req ListCommentReq
	ctxUtils := &GinCtxUtils{Context: c}
	if !ctxUtils.GetReqObject(&req) {
		return
	}

	daoObj := model.NewModel()

	list, err := daoObj.ListCommentByBlogID(req.BlogID)
	if err != nil {
		log.Printf("%#v", err)
		ctxUtils.ReplyFailServer()
		return
	}

	if len(list) == 0 {
		ctxUtils.ReplyOk(&ListCommentResp{List: make([]*model.Comment, 0)})
		return
	}

	for _, item := range list {
		if item.SubCount != 0 {
			subList, err := daoObj.ListCommentByParentID(item.ID)
			if err != nil {
				log.Printf("%#v", err)
				ctxUtils.ReplyFailServer()
				return
			}
			item.Children = subList
		}
	}

	ctxUtils.ReplyOk(&ListCommentResp{List: list})
}
