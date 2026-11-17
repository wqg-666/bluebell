package models

import "time"

// Post 内存对齐概念，相同类型的数据放在一起,缩小在内存中占的大小
type Post struct {
	//string,可以使后端传很大的数给前端，自动解析为string,也可以使前端传给后端的string转为int
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id" `
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AouthorName      string             `json:"aouthor_name" `
	VoteNum          int64              `json:"vote_num" `
	*Post                               //嵌入帖子结构体
	*CommunityDetail `json:"community"` //嵌入社区帖子结构体
	//`json:"community"`额外套一层 community，否则把他的数据提到外层
}

//{
//    "code": 1000,
//    "msg": "success",
//    "data": {
//        "aouthor_name": "wqg",
//        "id": 14361167782744064,
//        "author_id": 14361026703134720,
//        "community_id": 1,
//        "status": 0,
//        "title": "轻松熊",
//        "content": "轻松小熊0v0",
//        "create_time": "2026-11-09T23:06:09Z",
//        "community": {
//            "id": 1,
//            "name": "Go",
//            "introduction": "Golang",
//            "create_time": "2016-11-01T08:10:10Z"
//        }
//    }
//}
