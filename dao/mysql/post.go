package mysql

import (
	"bluebell/models"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	err = db.Create(p).Error
	return
}

// GetPostById 根据id查询单个数据帖子
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)

	err = db.Select("post_id", "title", "content", "author_id", "community_id", "create_time").Where("post_id=?", pid).Error
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	//按照分数排序
	posts = make([]*models.Post, 0, 2) //初始化一个切片，0是长度，2是容量
	//page：当前要查询第几页。size：每页显示多少条数据。
	//LIMIT 起始位置, 查询数量
	err = db.Model(&models.Post{}).
		Select("post_id", "title", "content", "author_id", "community_id", "create_time").
		Order("create_time DESC").
		Offset(int((page - 1) * size)). //这里第一个参数是跳过多少条数据
		Limit(int(size)).               //第二个参数是查询多少条数据。
		Find(&posts).Error
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, len(ids))
	if len(ids) == 0 {
		return postList, nil
	}
	err = db.Model(&models.Post{}).
		Select("post_id", "title", "content", "author_id", "community_id", "create_time").
		Where("post_id IN ?", ids).
		Order("FIND_IN_SET(post_id, ?)").
		Find(&postList, strings.Join(ids, ",")).Error
	//Find() 的第二个参数提供这个占位符对应的值。
	//FIND_IN_SET(?)根据指定的 ID 顺序，对查询结果进行排序。FIND_IN_SET(要查找的值, 字符串列表)
	//假设 Redis 中取出来的帖子 ID 顺序是：
	//ids := []string{"103", "101", "105"}
	//你希望 MySQL 查询出来的帖子也按照这个顺序排列，而不是按照数据库默认顺序排列。
	//strings.Join(ids, ",") 会把： []string{"103", "101", "105"} 转换为：103,101,105
	return
}
