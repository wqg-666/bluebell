package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 处理注册请求的参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" `
	Password   string `json:"password" binding:"required" `
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 处理登录请求的参数
type ParamLogin struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required,eqfield=Password"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID 从请求中获取当前的用户，request.getCurrentUserID
	PostID    string `json:"post_id" binding:"required"`              //帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=1 -1 0"` //赞成票（1）反对票（-1）取消投票（0）
	//`json:"string"告诉 JSON 序列化/反序列化时，把数字按照字符串的形式处理。JSON传参时传string类型
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	//这个example可以显示在swagger生成的页面上
	Page        int64  `json:"page" form:"page" example:"score"` //排序依据
	Size        int64  `json:"size" form:"size" `
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空
}
type ParamCommunityPostList struct {
	CommunityID int64 `json:"community_id" form:"community_id"`
	ParamPostList
}
