package models

// User 结构体用来表示数据库中的用户记录：
// 用户注册时，[logic/user.go]会把请求参数里的用户名和密码整理成一个 User 实例，再交给数据库层保存：
type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string
}
