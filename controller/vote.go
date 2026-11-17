package controller

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVoteController godoc
// @Summary 给帖子投票
// @Description 对帖子赞成、反对或取消投票。direction 字段用字符串传递：1、-1 或 0
// @Tags 投票
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ParamVoteData true "投票信息"
// @Success 200 {object} SwaggerResponse "code 为 1000 表示成功，其他值表示业务错误"
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamVoteData)
	//c.ShouldBindJSON(&p)：从前端发送的 HTTP 请求中读取 JSON 数据，将 JSON 数据解析并绑定到 Go 结构体中。
	if err := c.ShouldBindJSON(p); err != nil {
		//errs    类型断言成功后，得到转换后的错误集合
		errs, ok := err.(validator.ValidationErrors)
		//类型断言.判断 err 里面存储的错误，是否属于 validator.ValidationErrors 类型。如果是，就把它转换成这个类型。
		//validator.ValidationErrors
		//这是 validator 库定义的一个类型，用来表示参数校验过程中产生的错误集合。
		//例如：
		//- PostID 不符合要求。
		//- Action 不符合要求。
		if !ok { //参数校验错误
			ResponseError(c, CodeInvalidParam)
			return
		}
		//翻译并除掉错误提示中的结构体标识
		errDate := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errDate)
		return
	}
	//获取当前请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		if errors.Is(err, redis.ErrVoteRepeated) || errors.Is(err, redis.ErrVoteTimeExpire) {
			ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
			return
		}
		if errors.Is(err, redis.ErrPostNotFound) {
			ResponseErrorWithMsg(c, CodeInvalidParam, "帖子不存在或帖子发布时间没有写入 Redis")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucess(c, nil)
}
