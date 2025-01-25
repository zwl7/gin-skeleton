package models

import (
	"errors"
	"gin-skeleton/database/mysql"
	"gorm.io/gorm"
)

// 默认表名为struct名的复数 对应users
type User struct {
	BaseModel
	Name     string
	Password string
	Age      int64
	IsDel    byte
}

// CreateUser 创建user
func (that *User) CreateUser(user *User) (err error) {
	err = mysql.DB.Create(user).Error
	return
}

func (that *User) CreateUsers(users []*User) (err error) {
	//users := []*User{
	//	User{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
	//	User{Name: "Jackson", Age: 19, Birthday: time.Now()},
	//}
	err = mysql.DB.Create(&users).Error
	return
}

func (that *User) Delete(id int64) (err error) {
	err = mysql.DB.Delete(&User{}, id).Error
	return
}

func (that *User) GetOneById(id int64, user *User) (bool, error) {
	//db.First(&user, 10)
	// SELECT * FROM users WHERE id = 10;

	if _err := mysql.DB.First(user, id).Error; _err != nil {

		if errors.Is(_err, gorm.ErrRecordNotFound) {
			//只是没查询到记录而已，不是啥致命的错误
			return false, nil
		}
		//其他错误
		return false, _err
	}

	//没有错误，返回true
	return true, nil
}

func (that *User) GetOneByName(name string, user *User) (bool, error) {

	if _err := mysql.DB.Where("name = ?", name).First(user).Error; _err != nil {

		if errors.Is(_err, gorm.ErrRecordNotFound) {
			//只是没查询到记录而已，不是啥致命的错误
			return false, nil
		}
		//其他错误
		return false, _err
	}

	//没有错误，返回true
	return true, nil
}

func (that *User) GetOneByIds(ids []int64, users []*User) (err error) {
	//db.Find(&users, []int{1,2,3})
	// SELECT * FROM users WHERE id IN (1,2,3);
	err = mysql.DB.Find(&users, ids).Error
	return
}

func (that *User) GetOne(user *User) (err error) {
	err = mysql.DB.Select("Name", "Age").Where("age", 20).Find(user).Error
	return
}
