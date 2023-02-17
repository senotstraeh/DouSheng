package controller

import (
	"DouSheng/Dao"
	"DouSheng/Service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction 点赞或者取消赞操作;
func FavoriteAction(c *gin.Context) {
	UserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(UserId, 10, 64)
	VideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(VideoId, 10, 64)
	ActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(ActionType, 10, 64)
	like := Dao.Like{}

	result := Dao.DB.Where("VideoId = ?", videoId).First(&like)
	if result.RowsAffected == 0 {
		if actionType == 2 {
			c.JSON(http.StatusOK,
				Dao.CommonResponse{StatusCode: 1,
					StatusMsg: "cancel like fail"})
			return
		}
		like.UserId = userId
		like.VideoId = videoId
		like.Cancel = 0
		err := Service.InsertLike(like)
		if err != nil {
			c.JSON(http.StatusOK,
				Dao.CommonResponse{StatusCode: 1,
					StatusMsg: "action like fail"})
			return
		}
		c.JSON(http.StatusOK,
			Dao.CommonResponse{StatusCode: 0,
				StatusMsg: "action like success"})
		return

	}

	err := Service.UpdateLike(userId, videoId, actionType)
	if err != nil {
		c.JSON(http.StatusOK,
			Dao.CommonResponse{StatusCode: 1,
				StatusMsg: "action like fail"})
		return
	}
	c.JSON(http.StatusOK,
		Dao.CommonResponse{StatusCode: 0,
			StatusMsg: "action like success"})
	return

}

// GetFavouriteList 获取点赞列表;
func GetFavouriteList(c *gin.Context) {
	Userid := c.Query("user_id")

	userId, _ := strconv.ParseInt(Userid, 10, 64)

	videos, err := Service.GetLikeVideoIdList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Dao.GetFavouriteListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "get favouriteList fail "},
		})
		return
	}

	c.JSON(http.StatusOK, Dao.GetFavouriteListResponse{
		Dao.CommonResponse{StatusCode: 0,
			StatusMsg: "get favouriteList success "},
		videos,
	})
}
