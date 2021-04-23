package model

import (
	"ginVue/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Category Category `gorm:"foreignkey:Cid"`
	Title    string   `gorm:"type:varchar(100);not null " json:"title"`
	Cid      int      `gorm:"type:int;not null  json:"cid"`
	Desc     string   `gorm:"type:varchar(200)" json:"desc"`
	Content  string   `gorm:"type:longtext" json:"content"`
	Img      string   `gorm:"type:varchar(100)" json:"img"`
}

// 这里不需要对文章是否存在做验证 文章名可以重名
// 查询文章是否存在  传入文章名字 返回结果状态码
//func CheckArticle(title string) (code int) {
//	var art Article
//	// 在Article的Article表里面查找title=title的记录 此时将title赋值给art.title
//	db.Select("id").Where("title = ?", title).First(&art)
//	// 如果文章对象的id大于0，说明文章表中有title=title的记录
//	if art.ID > 0 {
//		return errmsg.ERROR_ART_NOT_EXIST
//	}
//	// 如果没有记录就存在 就可以创建用户
//	return errmsg.SUCCESS
//}

// 新增文章  传入文章数据 返回结果状态码
// 需要将文章信息写入数据库 所以传入的时候传入分类对象
func CreateArticle(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 如果创建成功就返回成功
	return errmsg.SUCCESS

}

// 查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err := db.Preload("Category").Where("id = ?", id).First(&art).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCESS
}

// 查询分类下的所有文章（多个数据需要分页）
func GetCategoryArticle(id int, pageSize int, pageNum int) ([]Article, int, int) {
	var categoryArticleList []Article
	var total int
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid =?", id).Find(&categoryArticleList).Count(&total).Error
	// 这条数据库语句需要传入在那个表里面找记录
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return categoryArticleList, errmsg.SUCCESS, total

}

// 查询文章列表
func GetArticle(title string, pageSize int, pageNum int) ([]Article, int, int) {
	var arts []Article
	var total int

	if title == "" {
		db.Model(&arts).Count(&total)
		db.Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&arts)
		// 这条数据库语句需要传入在那个表里面找记录
		// gorm数据库做了两件事 一是传入的是事先定义的表的对象  二是将查找到的记录赋值给这个对象 使这个对象有值 这里就存在arts切片里面
		//if err != nil && err != gorm.ErrRecordNotFound {
		//	return nil, errmsg.ERROR, 0
		//}


		return arts, errmsg.SUCCESS, total
	}

	db.Preload("Category").Where("title LIKE ?", title+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&arts).Count(&total)
	//db.Model(&arts)
	// 这条数据库语句需要传入在那个表里面找记录
	// gorm数据库做了两件事 一是传入的是事先定义的表的对象  二是将查找到的记录赋值给这个对象 使这个对象有值 这里就存在arts切片里面
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return nil, errmsg.ERROR, 0
	//}
	return arts, errmsg.SUCCESS, total
	// 第一个业务场景：用户不传文章名title 查询所有文章
	// 第二个业务场景：用户查询根据文章名查询文章 用户传文章名title
}

// 删除文章
func DeleteArticle(id int) int {
	var art Article
	err := db.Where("id = ?", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

// 编辑文章
func EddArticle(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid

	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err := db.Model(&art).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
