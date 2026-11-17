package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id,title,content,author_id,community_id)
               values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id查询单个数据帖子
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	//按照分数排序
	sqlStr := `select
    post_id,title,content,author_id,community_id,create_time
    from post
    order by  create_time
    DESC
    limit ?,?
    `
	posts = make([]*models.Post, 0, 2) //初始化一个切片，0是长度，2是容量
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	//page：当前要查询第几页。size：每页显示多少条数据。
	//LIMIT 起始位置, 查询数量
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)`
	//FIND_IN_SET(?)根据指定的 ID 顺序，对查询结果进行排序。FIND_IN_SET(要查找的值, 字符串列表)
	//假设 Redis 中取出来的帖子 ID 顺序是：
	//ids := []string{"103", "101", "105"}
	//你希望 MySQL 查询出来的帖子也按照这个顺序排列，而不是按照数据库默认顺序排列。
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	//strings.Join(ids, ","),一个字符串拼接函数。strings.Join(字符串切片, 分隔符)
	//query	展开后的 SQL 语句
	//args	SQL 中占位符对应的参数
	//err	是否发生错误
	//args 大致是：
	//[]interface{}{
	//    "103",
	//    "101",
	//    "105",
	//    "103,101,105",
	//}
	//如果 sqlx.In() 展开 SQL 时发生错误，就直接返回。
	if err != nil {
		return nil, err
	}
	//重新绑定占位符,它可以让代码适配不同数据库。
	//sqlx.In() 默认生成的是 ? 占位符。
	//而 db.Rebind() 会根据当前数据库驱动，将占位符转换成对应的格式。
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
