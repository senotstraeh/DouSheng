package Service

import (
	"DouSheng/Config"
	"DouSheng/Dao"
	//"DouSheng/controller"
	"errors"
	"sync"
	"time"
)

// Save 保存视频记录
func Save(video Dao.TableVideo) error {

	result := Dao.DB.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CommentIdList 根据视频id获取评论id 列表
func CommentIdList(videoId int64) ([]Dao.Comment, error) {
	var commentList []Dao.Comment
	err := Dao.DB.Model(Dao.Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentList).Error
	if err != nil {

		return nil, errors.New("query commentlist fail")
	}
	return commentList, nil
}

func GetVideoListByVideoId(newvideo *Dao.Video, likeVideoId int64, wg *sync.WaitGroup, errrtn error) {
	var tablevideo Dao.TableVideo
	defer wg.Done()
	err := Dao.DB.Model(Dao.TableVideo{}).Where("Id = ?", likeVideoId).First(&tablevideo)
	if err != nil {
		errrtn = errors.New("get favorite video fail")
		return
	}
	newvideo.TableVideo = tablevideo

	author, err2 := GetTableUserById(tablevideo.AuthorId)
	if err2 != nil {
		errrtn = errors.New("get userinfo fail")
		return
	}
	newvideo.Author.Id = author.Id
	newvideo.Author.Name = author.Name
	res := int64(-1)
	res = Dao.DB.Model(Dao.Comment{}).Where("VideoId = ?", likeVideoId).RowsAffected
	if res < 0 {
		errrtn = errors.New("get comment count fail")
		return
	}
	newvideo.CommentCount = res

	res = -1
	res = Dao.DB.Model(Dao.Like{}).Where("VideoId = ?", likeVideoId).RowsAffected
	if res < 0 {
		errrtn = errors.New("get favorite count fail")
		return
	}
	newvideo.FavoriteCount = res
	newvideo.IsFavorite = true
	errrtn = nil
	return
}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]Dao.TableVideo, error) {
	videos := make([]Dao.TableVideo, Config.VideoCount)
	result := Dao.DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(Config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

func GetFeedVideoList(tablevideos []Dao.TableVideo) ([]Dao.Video, time.Time, error) {
	videos := make([]Dao.Video, 0, Config.VideoCount)
	wg := &sync.WaitGroup{}
	errlist := make([]error, Config.VideoCount, Config.VideoCount)
	for i := 0; i < Config.VideoCount; i++ {
		wg.Add(1)
		go GetVideoListByVideoId(&videos[i], tablevideos[i].Id, wg, errlist[i])
	}

	defer wg.Done()
	for i := 0; i < Config.VideoCount; i++ {
		if errlist[i] != nil {
			return videos, time.Now(), errors.New("get video fail")
		}
	}
	return videos, videos[Config.VideoCount-1].PublishTime, nil
}

// GetVideosByAuthorId
// 根据作者的id来查询对应数据库数据，并TableVideo返回切片
func GetVideosByAuthorId(authorId int64) ([]Dao.TableVideo, error) {

	var data []Dao.TableVideo

	result := Dao.DB.Where(&Dao.TableVideo{AuthorId: authorId}).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
