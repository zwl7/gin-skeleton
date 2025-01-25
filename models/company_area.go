package models

import (
	"gin-skeleton/database/mysql"
)

// 默认表名为struct名的复数 对应users
type CompanyAreas struct {
	ID            int64 `gorm:"primarykey"`
	Name          string
	Pinyin        string
	Des           string
	CompanyId     int64
	Operator      int64
	CompanyAreaId int64
	sort          int64
	Lat           float64
	Lng           float64
}

func (that *CompanyAreas) GetAll(users *[]CompanyAreas) (err error) {
	//db.Find(&users, []int{1,2,3})

	// SELECT * FROM users WHERE name <> 'jinzhu';
	err = mysql.DB.Find(users).Error
	return
}
