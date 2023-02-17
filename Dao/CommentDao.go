package Dao

import (
	"time"
)

// Comment
// 评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id          int64     //评论id
	UserId      int64     //评论用户id
	VideoId     int64     //视频id
	CommentText string    //评论内容
	CreateDate  time.Time //评论发布的日期mm-dd
	Cancel      int32     //取消评论为1，发布评论为0
}

type CommentInfo struct {
	Id         int64
	User       User
	Content    string
	Creat_data string //"mm-dd"
}
type CommentActionResponse struct {
	CommonResponse
	CommentInfo
}

type CommentListResponse struct {
	CommonResponse
	Comment_List []CommentInfo
}
