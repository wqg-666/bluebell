package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	//保存到redis中
	err = redis.CreatePost(p.ID, p.CommunityID)
	//3.返回
	return
}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合我们接口想用的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id查询社区详情信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// 与列表页使用相同的 Redis 统计，返回真实赞成票数。
	voteData, err := redis.GetPostVoteData([]string{strconv.FormatInt(pid, 10)})
	if err != nil {
		return nil, err
	}
	//把查到的接口数据拼接
	data = &models.ApiPostDetail{
		VoteNum:         voteData[0],
		AouthorName:     user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AouthorName:     user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Error("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostlist2", zap.Any("ids", ids))
	// 3.根据id去mysql数据库查询帖子详细信息
	//返回的顺序按照给定的id顺序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostlist2", zap.Any("posts", posts))
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AouthorName:     user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 2.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Error("GetCommunityPostIDsInOrder", zap.Any("ids", ids))
		return
	}
	// 3.根据id去mysql数据库查询帖子详细信息
	//返回的顺序按照给定的id顺序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostlist2", zap.Any("posts", posts))
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AouthorName:     user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListNew 将两个查询逻辑合二为一的函数
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		//查所有
		data, err = GetPostList2(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
