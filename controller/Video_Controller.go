package controller

import (
	"DouSheng/Config"
	"DouSheng/Dao"
	"DouSheng/Service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Dao.CommonResponse
	VideoList []Dao.Video `json:"video_list"`
	NextTime  int64       `json:"next_time"`
}

type VideoListResponse struct {
	Dao.CommonResponse
	VideoList []Dao.Video `json:"video_list"`
}

// 官方feed
var DemoUser = Dao.User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
var DemoVideos = []Dao.Video{
	{
		TableVideo: Dao.TableVideo{
			Id:       1,
			PlayUrl:  "https://www.w3schools.com/html/movie.mp4",
			CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		},
		Author: DemoUser,

		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		CommonResponse: Dao.CommonResponse{StatusCode: 0},
		VideoList:      DemoVideos,
		NextTime:       time.Now().Unix(),
	})
}

/*
// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	var lastTime = time.Now()
	if inputTime != "0" {
		inputTimeInt, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(inputTimeInt, 0)
	}
	//userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	//有登录状态，验证token
	//token, ok := c.GetQuery("token")
	tablevideos, err := Service.GetVideosByLastTime(lastTime)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	videos, nexttime, err := Service.GetFeedVideoList(tablevideos)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		CommonResponse: Dao.CommonResponse{StatusCode: 0, StatusMsg: "获取视频流成功"},
		VideoList:      videos,
		NextTime:       nexttime.Unix(),
	})
}*/

// Publish /publish/action/
// 视频文件先存在本机
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	title := c.PostForm("title")

	//videoName := Service.NewFileName()
	if err != nil {

		c.JSON(http.StatusOK, Dao.CommonResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	dst := fmt.Sprintf("./tmp/%s", data.Filename)
	err = c.SaveUploadedFile(data, dst)
	var video Dao.TableVideo
	video.PublishTime = time.Now()
	video.PlayUrl = Config.PlayUrl
	video.CoverUrl = Config.CoverUrl
	video.AuthorId = userId
	video.Title = title
	err = Service.Save(video)
	if err != nil {
		log.Printf("方法videoService.Publish(data, userId) 失败：%v", err)
		c.JSON(http.StatusOK, Dao.CommonResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Printf("方法videoService.Publish(data, userId) 成功")

	c.JSON(http.StatusOK, Dao.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	user, err := Dao.GetTableUserById(userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "获取用户失败"},
		})
		return
	}

	//依据用户id查询所有的视频，获取视频列表
	tablevideos, err := Service.GetVideosByAuthorId(userId)
	if err != nil {

		c.JSON(http.StatusOK, VideoListResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}
	videolist := make([]Dao.Video, 0, len(tablevideos))

	for _, video := range tablevideos {
		var videotemp Dao.Video
		videotemp.Author.Id = user.Id
		videotemp.Author.Name = user.Name
		videotemp.TableVideo = video
		videolist = append(videolist, videotemp)
	}

	log.Printf("调用videoService.List(%v)成功", userId)
	c.JSON(http.StatusOK, VideoListResponse{
		CommonResponse: Dao.CommonResponse{StatusCode: 0},
		VideoList:      videolist,
	})
}
