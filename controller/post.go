package controller

import (
	"bluebell/logic"
	"bluebell/models"
	settings "bluebell/setting"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler godoc
// @Summary 创建帖子
// @Description 使用当前登录用户创建帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Post true "帖子信息（community_id、title、content 必填）"
// @Success 200 {object} SwaggerResponse
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	//c.ShouldBindJSON() //validator -->binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSucess(c, nil)
}

// GetPostDetailHandler godoc
// @Summary 获取帖子详情
// @Description 根据帖子 ID 获取帖子、作者和社区详情
// @Tags 帖子
// @Produce json
// @Security BearerAuth
// @Param id path int true "帖子 ID"
// @Success 200 {object} SwaggerPostDetailResponse
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（从URL中帖子的id）GET("/post/:id")
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64) //10进制，int64
	if err != nil {
		zap.L().Error("get post detail with invalid param, failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSucess(c, data)
}

//query	参数放在 URL 的查询字符串

// GetPostListHandler godoc
// @Summary 获取帖子列表
// @Description 按数据库分页获取帖子列表
// @Tags 帖子
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认 1"
// @Param size query int false "每页数量，默认 10"
// @Success 200 {object} SwaggerPostListResponse
// @Router /posts/ [get]
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := GetPageInfo(c)
	//1.获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	//2.返回响应
	ResponseSucess(c, data)
}

// GetPostListHandler2 godoc
// @Summary 获取排序帖子列表
// @Description 按时间或分数排序分页查询；可选 community_id 仅查询指定社区
// @Tags 帖子
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认取配置 post_list.page"
// @Param size query int false "每页数量，默认取配置 post_list.size"
// @Param order query string false "排序方式：time 或 score，默认取配置 post_list.order" Enums(time, score)
// @Param community_id query int false "社区 ID；不传时查询所有社区"
// @Success 200 {object} SwaggerPostListResponse
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	//前端的GET请求(query string)参数：/api/v1/posts2?page=1&size=2&order=time
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  settings.Conf.PostListConfig.Page,
		Size:  settings.Conf.PostListConfig.Size,
		Order: settings.Conf.PostListConfig.Order,
	}
	//c.ShouldBind()根据请求的数据类型选择响应的方法去获取数据
	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才用这个方法取数据
	//ShouldBindQuery,将 URL 中的 Query 参数自动绑定到 Go 结构体中，并进行参数校验。
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
	}
	data, err := logic.GetPostListNew(p) //更新：合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	//2.返回响应
	ResponseSucess(c, data)
}

// 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	//前端的GET请求(query string)参数：/api/v1/posts2?page=1&size=2&order=time
//	//初始化结构体时指定初始参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: models.ParamPostList{
//			Page:  viper.GetInt64("post_list.page"),
//			Size:  viper.GetInt64("post_list.size"),
//			Order: viper.GetString("post_list.order"),
//		},
//	}
//	//c.ShouldBind()根据请求的数据类型选择响应的方法去获取数据
//	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才用这个方法取数据
//	//ShouldBindQuery,将 URL 中的 Query 参数自动绑定到 Go 结构体中，并进行参数校验。
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetPostCommunityList with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//	}
//
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//	}
//	//2.返回响应
//	ResponseSucess(c, data)
//}
