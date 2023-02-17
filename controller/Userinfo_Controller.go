package controller

import (
	"DouSheng/Dao"
	"DouSheng/Service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Register POST douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	u := Service.GetTableUserByUsername(username)
	if username == u.Name {
		c.JSON(http.StatusOK, Dao.UserLoginResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := Dao.TableUser{
			Name:     username,
			Password: Service.EnCoder(password),
		}
		if Service.InsertTableUser(&newUser) != true {
			println("Insert Data Fail")
		}
		u := Service.GetTableUserByUsername(username)
		token := Service.GenerateToken(u)
		log.Println("注册返回的id: ", u.Id)
		c.JSON(http.StatusOK, Dao.UserLoginResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 0},
			UserId:         u.Id,
			Token:          token,
		})
	}
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := Service.EnCoder(password)
	println(encoderPassword)

	u := Service.GetTableUserByUsername(username)

	if encoderPassword == u.Password {
		token := Service.GenerateToken(u)
		c.JSON(http.StatusOK, Dao.UserLoginResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 0},
			UserId:         u.Id,
			Token:          token,
		})
	} else {
		c.JSON(http.StatusOK, Dao.UserLoginResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "Username or Password Error"},
		})
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)

	u, err := Service.GetTableUserById(id)
	if err != nil {
		c.JSON(http.StatusOK, Dao.UserInfoResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		c.JSON(http.StatusOK, Dao.UserInfoResponse{
			CommonResponse: Dao.CommonResponse{StatusCode: 0},
			Userinfo:       Dao.User{Id: u.Id, Name: u.Name},
		})
	}
}
