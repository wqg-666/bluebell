package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	if err := db.Model(&models.User{}).Where("username=?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	if viper.GetString("auth.password_salt") == "" {
		return errors.New("auth.password_salt or BLUEBELL_PASSWORD_SALT must be configured")
	}
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	err = db.Create(&user).Error
	return
}

// Login 用户登录
func Login(user *models.User) (err error) {
	if viper.GetString("auth.password_salt") == "" {
		return errors.New("auth.password_salt or BLUEBELL_PASSWORD_SALT must be configured")
	}
	oPassword := user.Password // 用户传入的登录密码

	// 根据用户名查询用户
	err = db.Select("user_id", "username", "password").
		Where("username = ?", user.Username).
		First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}

	return nil
}

// encryptPassword 把密码转换成 MD5 十六进制字符串的函数
func encryptPassword(opassword string) string {
	h := md5.New()
	_, _ = h.Write([]byte(opassword + viper.GetString("auth.password_salt")))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUserById 根据id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)

	err = db.Select("user_id", "username").Where("user_id=", uid).First(&user).Error

	return
}
