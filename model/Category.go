package model

import (
	"ginVue/utils/errmsg"
	"github.com/jinzhu/gorm"
)

// 分类结构体
type Category struct {
	ID   uint   `gorm:"primary key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null " json:"name"`
}

// 查询用户是否存在  传入用户名字 返回结果状态码
func CheckCategory(name string) (code int) {
	var cate Category
	// 在category的category表里面查找username=name的记录 此时将name赋值给users.name
	db.Select("id").Where("name = ?", name).First(&cate)
	// 如果cate对象的id大于0，说明cate表中有name=name的记录
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	// 如果没有记录就存在 就可以创建用户
	return errmsg.SUCCESS
}

// 新增分类  传入用户数据 返回结果状态码
// 需要将分类信息写入数据库 所以传入的时候传入分类对象
func CreateCategory(data *Category) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 如果创建成功就返回成功
	return errmsg.SUCCESS

}

// 查询分类列表
func GetCategory(pageSize int, pageNum int) ([]Category, int) {
	var cates []Category
	var total int
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cates).Count(&total).Error
	// 这条数据库语句需要传入在那个表里面找记录
	// gorm数据库做了两件事 一是传入的是事先定义的表的对象  二是将查找到的记录赋值给这个对象 使这个对象有值 这里就存在user切片里面
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	return cates, total
}

// 查询单个分类
func GetCateInfo(id int) (Category, int) {
	var cate Category
	db.Where("id = ?", id).First(&cate)
	return cate, errmsg.SUCCESS
}

// 删除分类
func DeleteCategory(id int) int {
	var cate Category
	err := db.Where("id = ?", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

// 编辑分类
func EddCategory(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err := db.Model(&cate).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
