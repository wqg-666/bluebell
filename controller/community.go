package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//注释	作用	具体含义
//Summary	接口概要	简短介绍接口功能
//Description	接口描述	更详细地说明接口用途
//Tags	接口分类	把接口归类到社区分组
//Produce	响应格式	表示接口返回 JSON 数据
//Security	身份认证	表示接口需要 BearerAuth 认证
//Success	成功响应	描述请求成功时返回的数据
//Router	路由信息	接口路径和请求方法

// CommunityHandler godoc
// @Summary 获取社区列表
// @Description 获取所有社区的 ID 和名称
// @Tags 社区
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerCommunityListResponse
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	//查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端保存暴露给外面
		return
	}
	ResponseSucess(c, data)
}

//Param id path int true "社区 ID"
//这是用来描述接口需要接收的参数的。
//Param 参数名称 参数位置 参数类型 是否必填 参数描述
//path	参数位置	URL 路径中

//Swagger 参数位置	Gin 获取方式	示例
//path	c.Param()	/community/1
//query	c.Query()	/community?id=1
//header	c.GetHeader()	Authorization: xxx
//body	c.ShouldBindJSON()	JSON 请求体

// CommunityDetailHandler godoc
// @Summary 获取社区详情
// @Description 根据社区 ID 获取社区信息
// @Tags 社区
// @Produce json
// @Security BearerAuth
// @Param id path int true "社区 ID"
// @Success 200 {object} SwaggerCommunityDetailResponse
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端保存暴露给外面
		return
	}
	ResponseSucess(c, data)
}
