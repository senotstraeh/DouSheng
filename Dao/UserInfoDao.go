package Dao

import "log"

// TableUser 对应数据库User表结构的结构体
type TableUser struct {
	Id       int64
	Name     string
	Password string
}

// User 最终封装后,controller返回的User结构体
type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	//TotalFavorited int64  `json:"total_favorited,omitempty"`
	//FavoriteCount  int64  `json:"favorite_count,omitempty"`
}
type UserLoginResponse struct {
	CommonResponse
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	CommonResponse
	Userinfo User `json:"user"`
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) (TableUser, error) {
	tableUser := TableUser{}
	if err := DB.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetTableUserById 根据user_id获得TableUser对象
func GetTableUserById(id int64) (TableUser, error) {
	tableUser := TableUser{}
	if err := DB.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *TableUser) bool {
	if err := DB.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
