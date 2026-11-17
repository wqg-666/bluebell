package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//这些不是 Go 代码，而是 Swagger 专门识别的注释。
//godoc：表示这是该函数的文档注释。
//Accept 在这里表示接口接收的请求数据格式。
//Produce：响应的数据格式。
//Param request body models.ParamSignUp true "注册信息"这是用来描述接口需要接收什么参数的。
//Param 参数名称 参数位置 参数类型 是否必填 参数描述
//这里的request 是给这组请求参数起的一个名字。body 表示参数放在 HTTP 请求体中。
//Success 200 {object} SwaggerResponse "code 为 1000 表示成功，其他值表示业务错误" SwaggerResponse	响应结构体名称

// SignUpHandler godoc
// @Summary 用户注册
// @Description 创建一个新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.ParamSignUp true "注册信息"
// @Success 200 {object} SwaggerResponse "code 为 1000 表示成功，其他值表示业务错误"
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	//判断这个字段类型对不对，请求参数格式是json格式
	p := new(models.ParamSignUp) //会创建一个models.ParamSignUp，并让 p 保存它的指针，类型是 *models.ParamSignUp。
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误
		//记录日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//请求 JSON 格式错误时，可能是 JSON 解析错误。
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		//JSON 格式正确、但字段不符合 binding 校验规则时，通常是 validator.ValidationErrors
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2.业务处理,从controller走到logic/serve这一层，逻辑（服务）层，负责处理业务逻辑
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSucess(c, nil)
}

// LoginHandler godoc
// @Summary 用户登录
// @Description 使用用户名和密码登录并获取 JWT
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.ParamLogin true "登录信息"
// @Success 200 {object} SwaggerLoginResponse "成功时 data 中包含 token、user_id 和 user_name"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误
		//记录日志
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		//validator.ValidationErrors	参数校验失败
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//JSON 格式正确、但字段不符合 binding 校验规则时，通常是 validator.ValidationErrors
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	ResponseSucess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //前端显示的最大值1<<53-1  int64类型的最大值1<<63-1会有问题，所以这里转为字符串
		"user_name": user.Username,
		"token":     user.Token,
	})

}
