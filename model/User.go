package model

import (
	"encoding/base64"
	"ginVue/utils/errmsg"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
	"log"
)

// 这个文件有model层和对数据库操作的dao层

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:int;not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type:int；DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// 查询用户是否存在  传入用户名字 返回结果状态码
func CheckUser(name string) (code int) {
	var user User
	// 在users的user表里面查找username=name的记录 此时将name赋值给users.name
	db.Select("id").Where("username = ?", name).First(&user)
	// 如果users对象的id大于0，说明user表中有username=name的记录
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	// 如果没有记录就存在 就可以创建用户
	return errmsg.SUCCESS
}

// 更新用户
func CheckUpUser(id int,name string) (code int) {
	var user User
	// 在users的user表里面查找username=name的记录 此时将name赋值给users.name
	db.Select("id，username").Where("username = ?", name).First(&user)
	if user.ID == uint(id) {
		return errmsg.SUCCESS
	}
	// 如果users对象的id大于0，说明user表中有username=name的记录
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	// 如果没有记录就存在 就可以创建用户
	return errmsg.SUCCESS
}

// 新增用户  传入用户数据 返回结果状态码
// 需要将用户信息写入数据库 所以传入的时候传入用户对象
func CreateUser(data *User) int {
	data.Password = ScryptPassWord(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 如果创建成功就返回成功
	return errmsg.SUCCESS

}

// 查询单个用户 才能编辑它 根据id在user表里面查
func GetUser(id int) (User, int) {
	var user User
	err :=db.Where("id = ?",id).First(&user).Error
	if err != nil {
		return user,errmsg.ERROR_USER_NOT_EXIST
	}
	return user,errmsg.SUCCESS
}

// 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int) {
	var users []User
	var total int
	//var user User
	//db.Model(&user).Count(&total)
	if username == "" {
		 db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total)
		 return users,total
	}else {
	 db.Where("username LIKE ?", username+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total)
	}
	// 这条数据库语句需要传入在那个表里面找记录
	// gorm数据库做了两件事 一是传入的是事先定义的表的对象  二是将查找到的记录赋值给这个对象 使这个对象有值 这里就存在user切片里面

	//if  err != gorm.ErrRecordNotFound {
		//return nil, 0
	//}
	return users, total
}

// 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

// 编辑用户
func EddUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role

	err := db.Model(&user).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 密码加密 (传入string类型的password，加密后的string类型的密码)
func ScryptPassWord(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 23, 34, 56, 67, 78, 45, 65}
	// 生成密钥
	hashPassWord, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	// 密钥加密
	finalPassWord := base64.StdEncoding.EncodeToString(hashPassWord)
	return finalPassWord
}

// 登陆验证 (根据传入的用户名获取到数据库密码，对密码进行加密后与数据库加密的密码进行比对）
func CheckLogin(username string, password string) int {
	var user User
	db.Where("username = ?", username).First(&user)
	// 1.判断用户是否存在
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	// 2.判断用户的密码是否正确
	if ScryptPassWord(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	// 3.判断用户是否为管理员
	if user.Role != 2 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCESS
}
