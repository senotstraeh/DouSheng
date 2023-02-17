package Service

import (
	"DouSheng/Dao"
	"errors"
	"log"
	"sync"
)

// UpdateLike 根据userId，videoId,actionType点赞或者取消赞
func UpdateLike(userId int64, videoId int64, actionType int64) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := Dao.DB.Model(Dao.Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	//如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	//更新操作成功
	return nil
}

// InsertLike 插入点赞数据
func InsertLike(likeData Dao.Like) error {
	//创建点赞数据，默认为点赞，cancel为0，返回错误结果
	err := Dao.DB.Model(Dao.Like{}).Create(&likeData).Error
	//如果有错误结果，返回插入失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

// GetLikeVideoIdList 根据userId查询所属点赞全部videoId
func GetLikeVideoIdList(userId int64) ([]Dao.Video, error) {
	var likeVideoIdList []int64
	var likeVideoList []Dao.Video
	err := Dao.DB.Model(Dao.Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": 0}).
		Pluck("video_id", &likeVideoIdList).Error
	if err != nil {

		return likeVideoList, errors.New("get likeVideoIdList failed")

	}
	wg := &sync.WaitGroup{}
	errlist := make([]error, 0, len(likeVideoIdList))
	likeVideoList = make([]Dao.Video, 0, len(likeVideoIdList))
	for i := 0; i < len(likeVideoIdList); i++ {
		wg.Add(1)

		go GetVideoListByVideoId(&likeVideoList[i], likeVideoIdList[i], wg, errlist[i])

	}
	wg.Wait()
	for i := 0; i < len(likeVideoIdList); i++ {
		if errlist[i] != nil {
			return likeVideoList, errors.New("get likeVideoIdList failed")
		}

	}
	return likeVideoList, nil
}
