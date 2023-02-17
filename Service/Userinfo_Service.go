package Service

import (
	"DouSheng/Dao"
	"log"
)

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) Dao.TableUser {
	tableUser, err := Dao.GetTableUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *Dao.TableUser) bool {
	flag := Dao.InsertTableUser(tableUser)
	if flag == false {
		log.Println("插入失败")
		return false
	}
	return true
}

// GetTableUserById 根据user_id获得TableUser对象
func GetTableUserById(id int64) (Dao.TableUser, error) {
	tableUser := Dao.TableUser{}
	if err := Dao.DB.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// EnCoder 密码加密待实现
func EnCoder(password string) string {
	return password
}

// GenerateToken 根据信息创建token,待实现
func GenerateToken(u Dao.TableUser) string {
	return string(u.Id) + " " + u.Name + u.Password
}
