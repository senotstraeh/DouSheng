package controller

import (
	"DouSheng/Dao"
	"DouSheng/Service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// CommentAction
// 发表 or 删除评论 comment/action/
func CommentAction(c *gin.Context) {
	log.Println("CommentController-Comment_Action: running") //函数已运行
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	log.Printf("err:%v", err)
	log.Printf("userId:%v", userId)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, Dao.CommentActionResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1,
				StatusMsg: "comment userId json invalid"},
		})
		log.Println("CommentController-Comment_Action: return comment userId json invalid") //函数返回userId无效
		return
	}
	//获取videoId
	videoid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, Dao.CommentActionResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1,
				StatusMsg: "comment videoid json invalid"},
		})
		log.Println("CommentController-Comment_Action: return comment userid json invalid") //函数返回userId无效
		return
	}
	//获取操作类型
	actiontype, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	//错误处理
	if err != nil || actiontype < 1 || actiontype > 2 {
		c.JSON(http.StatusOK, Dao.CommentActionResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1,
				StatusMsg: "comment actiontype json invalid"},
		})
		log.Println("CommentController-Comment_Action: return comment actiontype json invalid") //函数返回userId无效
		return
	}

	if actiontype == 2 {
		commentid_str := c.Query("comment_id")
		commentid, err := strconv.ParseInt(commentid_str, 10, 32)
		if err != nil {
			c.JSON(http.StatusOK, Dao.CommentActionResponse{
				CommonResponse: Dao.CommonResponse{StatusCode: 1,
					StatusMsg: "commentid  invalid"},
			})
			log.Println("CommentController-Comment_Action: return commentid json invalid") //函数返回commentid无效
			return
		}
		err = Service.DeleteComment(commentid)
		if err != nil {
			c.JSON(http.StatusOK, Dao.CommentActionResponse{
				CommonResponse: Dao.CommonResponse{StatusCode: 1,
					StatusMsg: "commentid  delete fail"},
			})
			log.Println("CommentController-Comment_Action: return commentid  delete fail") //函数返回删除评论无效
			return
		}
		c.JSON(http.StatusOK, Dao.CommentActionResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 0, StatusMsg: "comment delete success"},
		})
		return
	}
	//发表评论操作
	content := c.Query("comment_text")
	//针对内容可以做些敏感词屏蔽操作，目前还没做
	var toPublicComment Dao.Comment
	toPublicComment.UserId = userId
	toPublicComment.VideoId = videoid
	toPublicComment.CommentText = content
	toPublicComment.CreateDate = time.Now()
	toPublicComment.Cancel = 0
	//发表评论
	commentInfo, err := Service.PublicComment(toPublicComment)
	if err != nil {
		c.JSON(http.StatusOK, Dao.CommentActionResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1,
				StatusMsg: "comment public fail"},
		})
		log.Println("CommentController-Comment_Action: return comment public fail") //函数返回删除评论无效
		return
	}
	c.JSON(http.StatusOK, Dao.CommentActionResponse{
		CommonResponse: Dao.CommonResponse{StatusCode: 1,
			StatusMsg: "comment public success"},
		CommentInfo: commentInfo,
	})
	return
}

// 获取视频评论列表
func CommentList(c *gin.Context) {

	//获取videoId
	videoid := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Dao.CommentListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: -1,
				StatusMsg: " videoId json invalid"},
		})

		return
	}
	commentlist, err := Service.CommentIdList(videoId)
	if err != nil {
		c.JSON(http.StatusOK, Dao.CommentListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: -1,
				StatusMsg: " query Commentlist invalid"},
		})

		return
	}
	if len(commentlist) == 0 {
		c.JSON(http.StatusOK, Dao.CommentListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 0,
				StatusMsg: " empty comment"},
		})

		return
	}
	commentInfoList := make([]Dao.CommentInfo, len(commentlist))
	errlist := make([]error, len(commentlist))
	wg := &sync.WaitGroup{}
	for i := 0; i < len(commentInfoList); i++ {
		wg.Add(1)
		go Service.GetCommentInfoByCommentId(&commentInfoList[i], commentlist[i], wg, errlist[i])

	}
	wg.Wait()

	for i := 0; i < len(errlist); i++ {
		if errlist[i] != nil {
			c.JSON(http.StatusOK, Dao.CommentListResponse{
				CommonResponse: Dao.CommonResponse{StatusCode: 0,
					StatusMsg: " query comment fail"},
			})

			return
		}
	}

	//成功返回
	c.JSON(http.StatusOK, Dao.CommentListResponse{
		CommonResponse: Dao.CommonResponse{StatusCode: 0,
			StatusMsg: " query comment fail"},
		Comment_List: commentInfoList,
	})

	return

}
