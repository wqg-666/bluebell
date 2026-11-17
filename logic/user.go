package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) error {
	//1.判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	//`p.Username` 是请求 JSON 中 username 字段绑定得到的值。
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	//`p.Username` 是前端请求的 JSON 中 username 字段绑定得到的值。
	//这里的 &models.User{} 创建了一个 User 结构体，并返回它的指针。
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//Login(user),传递的user是个指针，mysql\user\login(*models.User)中db.Get(user, sqlStr, user.Username)。它会执行 SQL 查询，并将查询结果映射到 user 指向的结构体中。就能拿到user.UserID，
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT的token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
