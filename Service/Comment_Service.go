package Service

import (
	"DouSheng/Dao"
	"errors"
	"log"
	"sync"
)

// DeleteComment
// 3、删除评论，传入评论id
func DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	var comment Dao.Comment
	//先查询是否有此评论
	result := Dao.DB.Model(Dao.Comment{}).Where(map[string]interface{}{"id": id, "cancel": 0}).First(&comment)
	if result.RowsAffected == 0 { //查询到此评论数量为0则返回无此评论
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := Dao.DB.Model(Dao.Comment{}).Where("id = ?", id).Update("cancel", 1).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

func PublicComment(comment Dao.Comment) (Dao.CommentInfo, error) {
	err := Dao.DB.Model(Dao.Comment{}).Create(&comment).Error

	if err != nil {

		return Dao.CommentInfo{}, errors.New("create comment failed")
	}
	tableUser, err := Dao.GetTableUserById(comment.UserId)
	if err != nil {

		return Dao.CommentInfo{}, errors.New("query userinfo failed")
	}
	user := Dao.User{
		Id:            tableUser.Id,
		Name:          tableUser.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	var ccc Dao.CommentInfo
	ccc.Id = comment.Id
	ccc.User = user

	return Dao.CommentInfo{
		Id:         comment.Id,
		User:       user,
		Content:    comment.CommentText,
		Creat_data: comment.CreateDate.Format("2006-10-10 12:00:00.000"),
	}, nil
}

// 通过CommentId查询CommentInfo
func GetCommentInfoByCommentId(commentinfo *Dao.CommentInfo, commentdemo Dao.Comment, wg *sync.WaitGroup, err error) {
	//社交部分等待扩展
	user := Dao.User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	defer wg.Done()
	tableuser, err_get := GetTableUserById(commentdemo.UserId)
	if err_get != nil {
		err = err_get
		return
	}
	user.Id = tableuser.Id
	user.Name = tableuser.Name
	commentinfo.Id = commentdemo.Id
	commentinfo.Content = commentdemo.CommentText
	commentinfo.Creat_data = commentdemo.CreateDate.Format("2006-10-10 12:00:00.000")
	commentinfo.User = user
	return
}
